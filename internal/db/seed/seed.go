package seed

import (
	"database/sql"
	"embed"
	"vietio/internal/ads"
)

//go:embed images/*
var imagesFS embed.FS

func Run(dbConn *sql.DB, fileStorage ads.FileStorage) error {
    if err := runCategoriesSeed(dbConn); err != nil {
        return err
    }

    if err := runCitiesSeed(dbConn); err != nil {
        return err
    }

    if err := runUsersSeed(dbConn); err != nil {
        return err
    }

    if err := runAdsSeed(dbConn); err != nil {
        return err
    }

    if err := runFilesSeed(dbConn, fileStorage, imagesFS); err != nil {
        return err
    }

    return nil
}
