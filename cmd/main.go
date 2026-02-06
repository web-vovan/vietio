package main

import (
	"database/sql"
	"io"
	"log/slog"
	"os"
	"vietio/config"
	"vietio/internal/app"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
    logger, logFile := setupLogger()
	defer logFile.Close()

	slog.SetDefault(logger)

    config := config.Load()

    dbConn, err := sql.Open("pgx", config.Db.Dsn)
    if err != nil {
        logger.Error("db not connection", "err", err)
        os.Exit(1)
    }
    defer dbConn.Close()

    if err := dbConn.Ping(); err != nil {
        logger.Error("db ping failed", "err", err)
        os.Exit(1)
    }

    if config.SeedFlag {
        app.RunSeed(dbConn, config, logger)
		return
    }
        
	app.RunMigrations(dbConn, logger)

    app.RunHttpServer(dbConn, config, logger)
}

func setupLogger() (*slog.Logger, *os.File) {
	file, err := os.OpenFile(
		"app.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		panic(err)
	}

	multiWriter := io.MultiWriter(os.Stdout, file)

	handler := slog.NewJSONHandler(
		multiWriter,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	)

	logger := slog.New(handler)

	return logger, file
}