package seed

import (
	"database/sql"
)

func runCategoriesSeed(dbConn *sql.DB) error {
    _, err := dbConn.Exec(`
        INSERT INTO categories ("name", "order")
        SELECT *
        FROM (
            VALUES
                ('Барахолка', 1),
                ('Байки', 2),
                ('Жильё', 3),
                ('Услуги', 4),
                ('Разное', 5)
        ) AS v("name", "order")
        WHERE NOT EXISTS (SELECT 1 FROM categories);
    `)
    
    if err != nil {
        return err
    }

    return nil
}