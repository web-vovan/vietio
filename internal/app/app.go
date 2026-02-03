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
	"vietio/internal/middleware"
	"vietio/internal/storage"
	"vietio/internal/user"
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

	authValidator := auth.NewValidator()
	userRepository := user.NewRepository(dbConn)
	authService := auth.NewService(config, authValidator, userRepository)
	authHandler := auth.NewHandler(authService)

	// middleware
	authMiddleware := middleware.AuthJWT(authService)

	router := http.NewServeMux()

	// публичные роуты
	router.HandleFunc("GET /api/ads", adsHandler.GetAds)
	router.HandleFunc("POST /api/auth/login", authHandler.GetToken)

	// роуты с авторизацией
	router.Handle(
		"GET /api/ads/{uuid}",
		authMiddleware(http.HandlerFunc(adsHandler.GetAd)),
	)
	router.Handle(
		"POST /api/ads",
		authMiddleware(http.HandlerFunc(adsHandler.CreateAd)),
	)
	router.Handle(
		"PUT /api/ads/{uuid}", 
		authMiddleware(http.HandlerFunc(adsHandler.UpdateAd)),
	)
	router.Handle(
		"DELETE /api/ads/{uuid}", 
		authMiddleware(http.HandlerFunc(adsHandler.DeleteAd)),
	)

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
