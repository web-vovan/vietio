package seed

import (
	"database/sql"
)

func runCategoriesSeed(dbConn *sql.DB) error {
    _, err := dbConn.Exec(`
        INSERT INTO categories ("name", "order")
        VALUES
            ('Барахолка', 1),
            ('Байки', 2),
            ('Жильё', 3),
            ('Услуги', 4)
            ('Разное', 5)
    `)
    
    if err != nil {
        return err
    }

    return nil
}