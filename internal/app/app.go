package app

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"vietio/config"
	"vietio/internal/ads"
	"vietio/internal/auth"
	"vietio/internal/categories"
	"vietio/internal/db/seed"
	"vietio/internal/file"
	"vietio/internal/storage"
	"vietio/migrations"
)

func RunMigrations(dbConn *sql.DB) {
	if err := migrations.Up(dbConn); err != nil {
		log.Fatal("migration failed: ", err)
	}

	log.Println("успешная миграция БД")
}

func RunSeed(dbConn *sql.DB, config *config.Config) {
	if config.Env != "dev" {
		log.Fatal("сиды работают только в dev окружении")
	}

	if err := migrations.Reset(dbConn); err != nil {
		log.Fatal("reset failed: ", err)
	}

	log.Println("rollback всех таблиц")

	if err := migrations.Up(dbConn); err != nil {
		log.Fatal("migration failed: ", err)
	}

	log.Println("успешная миграция БД")

	if err := seed.Run(dbConn); err != nil {
		log.Fatal("seed failed: ", err)
	}

	log.Println("сиды успешно добавлены")
}

func RunHttpServer(dbConn *sql.DB, config *config.Config) {
	adsRepository := ads.NewRepository(dbConn)
	categoryRepository := categories.NewRepository(dbConn)
	fileRepository := file.NewFileRepository(dbConn)
	adValidator := ads.NewValidator(categoryRepository, adsRepository)

	var fileStorage ads.FileStorage
	var err error

	switch config.StorageType {
	case "local":
		fileStorage = storage.NewLocalStorage(config.Server.PublicUrl, "./uploads")
	case "s3":
		fileStorage, err = storage.NewS3Storage(
			context.Background(),
			config.S3Storage.Key,
			config.S3Storage.Secret,
			config.S3Storage.Bucket,
			config.S3Storage.PublicUrl,
		)
		if err != nil {
			panic("failed to init s3 storage: " + err.Error())
		}
	default:
		panic("неизвестный тип хранилища: " + config.StorageType)
	}

	adsService := ads.NewService(
		adsRepository,
		fileRepository,
		fileStorage,
		adValidator,
	)
	adsHandler := ads.NewHandler(adsService)

	authService := auth.NewService(config)
	authHandler := auth.NewHandler(authService)

	router := http.NewServeMux()
	router.HandleFunc("GET /api/ads", adsHandler.GetAds)
	router.HandleFunc("POST /api/ads", adsHandler.CreateAd)
	router.HandleFunc("GET /api/ads/{uuid}", adsHandler.GetAd)
	router.HandleFunc("PUT /api/ads/{uuid}", adsHandler.UpdateAd)
	router.HandleFunc("POST /api/auth/login", authHandler.GetToken)

	if config.Env == "dev" {
		router.HandleFunc("/api/test-init-data/{username}", authHandler.GetTestInitData)
	}

	// отдаем статику для локального хранилища
	if config.StorageType == "local" {
		router.Handle(
			"/uploads/",
			http.StripPrefix(
				"/uploads/",
				http.FileServer(http.Dir("./uploads")),
			),
		)
	}

	server := http.Server{
		Addr:    ":" + config.Server.HttpPort,
		Handler: router,
	}

	server.ListenAndServe()
}

