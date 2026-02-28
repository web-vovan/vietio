package seed

import (
	"database/sql"
	"embed"
	"log/slog"
	"vietio/internal/ads"
)

//go:embed images/*
var imagesFS embed.FS

func Run(dbConn *sql.DB, fileStorage ads.FileStorage, logger *slog.Logger) error {
    if err := runCategoriesSeed(dbConn); err != nil {
        return err
    }
    logger.Info("seed categories success")

    if err := runCitiesSeed(dbConn); err != nil {
        return err
    }

    logger.Info("seed cities success")

    if err := runUsersSeed(dbConn); err != nil {
        return err
    }
    logger.Info("seed users success")

    if err := runAdsSeed(dbConn); err != nil {
        return err
    }
    logger.Info("seed ads success")

    if err := runFilesSeed(dbConn, fileStorage, imagesFS); err != nil {
        return err
    }
    logger.Info("seed files success")

    return nil
}
