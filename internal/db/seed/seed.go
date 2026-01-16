package seed

import (
	"database/sql"
)

func Run(dbConn *sql.DB) error {
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

    if err := runFilesSeed(dbConn); err != nil {
        return err
    }

    return nil
}
