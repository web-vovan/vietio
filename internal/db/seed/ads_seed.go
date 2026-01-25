package seed

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

func runAdsSeed(dbConn *sql.DB) error {
    categoryIds, err := getAllIdsFromTable(dbConn, "categories")
	if err != nil {
		return err
	}

    userIds, err := getAllIdsFromTable(dbConn, "users")
	if err != nil {
		return err
	}

    districtList := [4]string{"юг", "центр", "север", "запад"}

    query := `
		INSERT INTO ads (
			user_id, category_id, city_id, title, description, price, currency, 
			district, status, expires_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	stmt, err := dbConn.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

    for i := 0; i < 150; i++ {
		userID := userIds[rand.Intn(len(userIds))]
		categoryID := categoryIds[rand.Intn(len(categoryIds))]
		cityID := 1
		title := gofakeit.ProductName()
		desc := gofakeit.Paragraph(1, 3, 10, " ")
		price := gofakeit.Number(10000, 50000000)
		currency := "VND"
		district := districtList[rand.Intn(len(districtList))]
		status := 1
		createdAt := gofakeit.DateRange(time.Now().AddDate(0, 0, -15), time.Now())
		updatedAt := createdAt
		expiresAt := createdAt.AddDate(0, 1, 0)

		_, err := stmt.Exec(
			userID,
			categoryID,
			cityID,
			title,
			desc,
			price,
			currency,
			district,
			status,
			expiresAt,
			createdAt,
			updatedAt,
		)
		if err != nil {
			return err
		}
	}

    return nil
}
