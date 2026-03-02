package app

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"vietio/config"
	"vietio/internal/ads"
	"vietio/internal/auth"
	"vietio/internal/categories"
	"vietio/internal/db/seed"
	"vietio/internal/file"
	"vietio/internal/middleware"
	"vietio/internal/storage"
	"vietio/internal/telegram"
	"vietio/internal/user"
	"vietio/internal/wishlist"
	"vietio/migrations"
)

func RunUpMigrations(dbConn *sql.DB, logger *slog.Logger) {
	if err := migrations.Up(dbConn); err != nil {
		logger.Error("migration failed", "err", err)
	}

	logger.Info("успешная миграция БД")
}

func RunDownMigrations(dbConn *sql.DB, logger *slog.Logger) {
	if err := migrations.Down(dbConn); err != nil {
		logger.Error("migration failed", "err", err)
	}

	logger.Info("успешный rollback миграции БД")
}

func RunSeed(dbConn *sql.DB, config *config.Config, logger *slog.Logger) {
	if config.Env != "dev" {
		logger.Error("сиды работают только в dev окружении")
		os.Exit(1)
	}

	fileStorage, err := getFileStorage(config, logger)
	if err != nil {
		os.Exit(1)
	}

	if err := migrations.Reset(dbConn); err != nil {
		logger.Error("reset failed", "err", err)
		os.Exit(1)
	}

	logger.Info("rollback всех таблиц")

	if err := migrations.Up(dbConn); err != nil {
		logger.Error("migration failed", "err", err)
		os.Exit(1)
	}

	logger.Info("успешная миграция БД")

	if err := seed.Run(dbConn, fileStorage, logger); err != nil {
		logger.Error("seed failed", "err", err)
		os.Exit(1)
	}

	logger.Info("сиды успешно добавлены")
}

func RunArchive(dbConn *sql.DB, config *config.Config, logger *slog.Logger) {
	adsRepository := ads.NewRepository(dbConn)
	categoryRepository := categories.NewRepository(dbConn)
	fileRepository := file.NewFileRepository(dbConn)
	userRepository := user.NewRepository(dbConn)
	wishlistRepository := wishlist.NewRepository(dbConn)
	adValidator := ads.NewValidator(categoryRepository, adsRepository)

	fileStorage, err := getFileStorage(config, logger)
	if err != nil {
		os.Exit(1)
	}

	adsService := ads.NewService(
		adsRepository,
		fileRepository,
		userRepository,
		wishlistRepository,
		fileStorage,
		adValidator,
	)

	err = adsService.ArchivingAds(context.Background())
	if err != nil {
		logger.Error("ошибка при отправке объявлений в архив", "err", err)
		os.Exit(1)
	}

	logger.Info("объявления успешно отправлены в архив")
}

func RunHttpServer(dbConn *sql.DB, config *config.Config, logger *slog.Logger) {
	adsRepository := ads.NewRepository(dbConn)
	categoryRepository := categories.NewRepository(dbConn)
	fileRepository := file.NewFileRepository(dbConn)
	userRepository := user.NewRepository(dbConn)
	wishlistRepository := wishlist.NewRepository(dbConn)
	adValidator := ads.NewValidator(categoryRepository, adsRepository)

	fileStorage, err := getFileStorage(config, logger)
	if err != nil {
		os.Exit(1)
	}

	adsService := ads.NewService(
		adsRepository,
		fileRepository,
		userRepository,
		wishlistRepository,
		fileStorage,
		adValidator,
	)
	adsHandler := ads.NewHandler(adsService, logger)

	authValidator := auth.NewValidator()
	authService := auth.NewService(config, authValidator, userRepository)
	authHandler := auth.NewHandler(authService)

	tgClient := telegram.NewClient(config.BotToken)
	telegramHandler := telegram.NewHandler(logger, tgClient)

	// middleware
	authMiddleware := middleware.AuthJWT(authService)

	router := http.NewServeMux()

	// публичные роуты
	router.HandleFunc("GET /api/ads", adsHandler.GetAds)
	router.HandleFunc("POST /api/auth/login", authHandler.GetToken)
	router.HandleFunc("POST /api/webhook", telegramHandler.Webhook)

	// роуты с авторизацией
	router.Handle(
		"GET /api/my",
		authMiddleware(http.HandlerFunc(adsHandler.GetMyAds)),
	)
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

	router.Handle(
		"POST /api/ads/{uuid}/sold",
		authMiddleware(http.HandlerFunc(adsHandler.MarkingSoldAd)),
	)

	router.Handle(
		"GET /api/my/sold",
		authMiddleware(http.HandlerFunc(adsHandler.GetMySoldAds)),
	)

	router.Handle(
		"POST /api/ads/{uuid}/favorite",
		authMiddleware(http.HandlerFunc(adsHandler.AddFavorite)),
	)

	router.Handle(
		"DELETE /api/ads/{uuid}/favorite",
		authMiddleware(http.HandlerFunc(adsHandler.DeleteFavorite)),
	)

	router.Handle(
		"GET /api/my/favorites",
		authMiddleware(http.HandlerFunc(adsHandler.GetMyFavoritesAds)),
	)

	// @todo убрать
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
		Handler: middleware.RecoverMiddleware(logger, router),
	}

	server.ListenAndServe()
}

func getFileStorage(config *config.Config, logger *slog.Logger) (ads.FileStorage, error) {
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
			logger.Error("failed to init s3 storage", "err", err)
			return nil, err
		}
	default:
		logger.Error("неизвестный тип хранилища", "storage", config.StorageType)
		return nil, err
	}

	return fileStorage, nil
}
