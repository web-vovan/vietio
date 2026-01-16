package seed

import (
	"database/sql"

    "github.com/brianvoe/gofakeit/v7"
)

func runUsersSeed(dbConn *sql.DB) error {
    for i := 0; i < 10; i++ {
		fakeTelegramID := int64(gofakeit.Number(100000000, 9999999999))
		fakeUsername := gofakeit.Username()

		query := `
			INSERT INTO users (telegram_id, username) 
			VALUES ($1, $2)
			ON CONFLICT (telegram_id) DO NOTHING`

		_, err := dbConn.Exec(query, fakeTelegramID, fakeUsername)

		if err != nil {
			return err
		}
	}

    return nil
}