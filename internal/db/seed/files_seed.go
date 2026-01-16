package seed

import (
	"database/sql"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
)

func runFilesSeed(dbConn *sql.DB) error {
    adsIds, err := getAllIdsFromTable(dbConn, "ads")
	if err != nil {
		return err
	}

    query := `
		INSERT INTO files (
			ad_id, path, "order"
		) VALUES ($1, $2, $3)
	`
	stmt, err := dbConn.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, id := range adsIds {
		for i := 0; i < 3; i++ {
			_, err := stmt.Exec(
				id,
				"/uploads/" + gofakeit.UUID() + ".jpg",
				i + 1,
			)

			if err != nil {
				return err
			}
		}
	}

	return nil
}
