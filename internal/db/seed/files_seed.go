package seed

import (
	"database/sql"
	"fmt"
	"math/rand"

	"github.com/brianvoe/gofakeit/v7"
)

func runFilesSeed(dbConn *sql.DB) error {
    adsUuids, err := getAllUuidFromTable(dbConn, "ads")
	if err != nil {
		return err
	}

    query := `
		INSERT INTO files (
			ad_uuid,
			path, 
			preview_path, 
			"order", 
			size, 
			preview_size, 
			mime, 
			preview_mime
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	stmt, err := dbConn.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, uuid := range adsUuids {
		for i := 0; i < 3; i++ {
			path := gofakeit.UUID()
			_, err := stmt.Exec(
				uuid,
				path + ".jpg",
				path + "_preview.jpg",
				i + 1,
				rand.Intn(1000),
				rand.Intn(1000),
				"image/jpg",
				"image/jpg",
			)

			if err != nil {
				return err
			}
		}
	}

	return nil
}
