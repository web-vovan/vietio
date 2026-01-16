package seed

import (
	"database/sql"
)

func runCategoriesSeed(dbConn *sql.DB) error {
    _, err := dbConn.Exec(`
        INSERT INTO categories ("name", "order")
        VALUES
            ('аренда байка', 1),
            ('аренда жилья', 2),
            ('личные вещи', 3),
            ('бесплатно', 4)
    `)

    if err != nil {
        return err
    }

    return nil
}