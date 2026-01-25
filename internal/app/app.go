package app

import (
	"database/sql"
	"log"
	"net/http"

	"vietio/config"
	"vietio/internal/ads"
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
	fileStorage := storage.NewLocalStorage(config.Server.PublicFilesBaseUrl, "./uploads")
	fileRepository := file.NewFileRepository(dbConn)

	adsService := ads.NewService(
		adsRepository,
		categoryRepository,
		fileStorage,
		fileRepository,
	)
	adsHandler := ads.NewHandler(adsService)

	router := http.NewServeMux()
	router.HandleFunc("GET /api/ads", adsHandler.GetAds)
	router.HandleFunc("POST /api/ads", adsHandler.CreateAd)
	router.HandleFunc("GET /api/ads/{uuid}", adsHandler.GetAd)

	// отдаем статику, в дальнейшем переедет в nginx
	router.Handle(
		"/uploads/",
		http.StripPrefix(
			"/uploads/",
			http.FileServer(http.Dir("./uploads")),
		),
	)

	server := http.Server{
		Addr:    ":" + config.Server.HttpPort,
		Handler: router,
	}

	server.ListenAndServe()
}
