package ads

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (repo *Repository) FindAds(ctx context.Context, params AdsListFilterParams) (AdsListRepository, error) {
	var result AdsListRepository

	var ads []AdsListItemRepository
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
			uuid,
			title,
            category_id,
            price,
            created_at,
			COALESCE(f.preview_path, '') as image,
            count(*) over() as total
		FROM ads
		LEFT JOIN LATERAL (
			SELECT preview_path
			FROM files
			WHERE files.ad_uuid = ads.uuid
			ORDER BY created_at ASC
			LIMIT 1
		) f ON true
        %s
        ORDER BY ads.%s %s
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
		var ad AdsListItemRepository
		if err := rows.Scan(
			&ad.Uuid,
			&ad.Title,
			&ad.CategoryId,
			&ad.Price,
			&ad.CreatedAt,
			&ad.Image,
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

func (repo *Repository) CreateAd(ctx context.Context, tx *sql.Tx, payload CreateAdRequestBody) (uuid.UUID, error) {
	var uuid uuid.UUID

	query := `
		INSERT INTO ads (
			title,
			description,
			category_id,
			price,
			user_id,
			city_id,
			currency,
			status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING uuid
	`

	err := tx.QueryRowContext(
		ctx,
		query,
		payload.Title,
		payload.Description,
		payload.CategoryId,
		payload.Price,
		1,
		1,
		"VDN",
		1,
	).Scan(&uuid)

	if err != nil {
		return uuid, err
	}
	return uuid, nil
}

func (repo *Repository) FindAdByUuid(ctx context.Context, uuid uuid.UUID) (AdModel, error) {
	var result AdModel

	query := `
        SELECT
			uuid,
			title,
            description,
            category_id,
            price,
            created_at
		FROM ads
		WHERE uuid = $1
		LIMIT 1
    `

	err := repo.db.QueryRowContext(ctx, query, uuid).Scan(
		&result.Uuid,
		&result.Title,
		&result.Description,
		&result.CategoryId,
		&result.Price,
		&result.CreatedAt,
	)
	if err != nil {
		return result, err
	}

	return result, nil
}
