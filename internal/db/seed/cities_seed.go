package seed

import (
	"database/sql"
)

func runCitiesSeed(dbConn *sql.DB) error {
    _, err := dbConn.Exec(`
        INSERT INTO cities ("name_vn", "name_rus")
        VALUES ('Nha Trang', 'Нячанг')
    `)

    if err != nil {
        return err
    }

    return nil
}