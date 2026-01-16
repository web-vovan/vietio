package ads

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (repo *Repository) FindAds(ctx context.Context, params AdsListFilterParams) (AdsListDB, error) {
	var result AdsListDB

	var ads []Ad
	var total int
	var conditions []string
	var args []any

	argsPos := 1

	if params.CategoryId != nil {
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", argsPos))
		args = append(args, *params.CategoryId)
		argsPos++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf(`
        SELECT
			id,
			title,
            description,
            category_id,
            price,
            created_at,
            count(*) over() as total
		FROM ads
        %s
        ORDER BY %s %s
		LIMIT %d OFFSET %d
    `,
		where,
		params.Sort,
		params.Order,
		params.Limit,
		params.Limit*(params.Page-1),
	)

	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var ad Ad
		if err := rows.Scan(
			&ad.ID,
			&ad.Title,
			&ad.Description,
			&ad.CategoryId,
			&ad.Price,
			&ad.CreatedAt,
			&total,
		); err != nil {
			return result, err
		}
		ads = append(ads, ad)
	}

	result.Items = ads
	result.Total = total

	return result, nil
}
