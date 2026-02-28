package seed

import (
	"database/sql"
)

func runCitiesSeed(dbConn *sql.DB) error {
    _, err := dbConn.Exec(`
        INSERT INTO cities ("name_vn", "name_rus")
        SELECT *
        FROM (
            VALUES ('Nha Trang', 'Нячанг')
        ) AS v("name_vn", "name_rus")
        WHERE NOT EXISTS (SELECT 1 FROM cities);
    `)

    if err != nil {
        return err
    }

    return nil
}