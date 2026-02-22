package ads

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"vietio/internal/authctx"

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

	if params.UserId != nil {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argsPos))
		args = append(args, *params.UserId)
		argsPos++
	}

	if params.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argsPos))
		args = append(args, params.Status)
		argsPos++
	} else {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argsPos))
		args = append(args, STATUS_ACTIVE)
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

	userId, err := authctx.GeUserIdFromContext(ctx)
	if err != nil {
		return uuid, err
	}

	query := `
		INSERT INTO ads (
			title,
			description,
			category_id,
			price,
			user_id,
			city_id,
			currency,
			status,
			expires_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_DATE + INTERVAL '1 month')
		RETURNING uuid
	`

	err = tx.QueryRowContext(
		ctx,
		query,
		payload.Title,
		payload.Description,
		payload.CategoryId,
		payload.Price,
		userId,
		1,
		"VDN",
		STATUS_ACTIVE,
	).Scan(&uuid)

	if err != nil {
		return uuid, err
	}
	return uuid, nil
}

func (repo *Repository) UpdateAd(ctx context.Context, tx *sql.Tx, ad AdModel) error {
	query := `
		UPDATE ads
		SET
			title = $1,
			description = $2,
			price = $3,
			category_id = $4,
			updated_at = now()
		WHERE 
			uuid = $5
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		ad.Title,
		ad.Description,
		ad.Price,
		ad.CategoryId,
		ad.Uuid,
	)
	if err != nil {
		return err
	}
	
	return nil
}

func (repo *Repository) FindAdByUuid(ctx context.Context, uuid uuid.UUID) (AdModel, error) {
	var result AdModel

	query := `
        SELECT
			uuid,
			title,
            description,
			user_id,
            category_id,
            price,
			status,
            created_at
		FROM ads
		WHERE uuid = $1
		LIMIT 1
    `

	err := repo.db.QueryRowContext(ctx, query, uuid).Scan(
		&result.Uuid,
		&result.Title,
		&result.Description,
		&result.UserId,
		&result.CategoryId,
		&result.Price,
		&result.Status,
		&result.CreatedAt,
	)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (repo *Repository) DeleteAdByUuidWithTx(ctx context.Context, tx *sql.Tx, uuid uuid.UUID) error {
	query := `
		DELETE FROM ads
		WHERE uuid = $1
	`
	
	_, err := tx.ExecContext(
		ctx,
		query,
		uuid,
	)
	if err != nil {
		return err
	}
	
	return nil
}

func (repo *Repository) ChangeStatusAdByUuidWithTx(ctx context.Context, tx *sql.Tx, status int, uuid uuid.UUID) error {
	query := `
		UPDATE ads
		SET
			status = $1,
			updated_at = now()
		WHERE 
			uuid = $2
	`
	
	_, err := tx.ExecContext(
		ctx,
		query,
		status,
		uuid,
	)
	
	if err != nil {
		return err
	}
	
	return nil
}

func (r *Repository) Exists(ctx context.Context, uuid uuid.UUID) (bool, error) {
    var result bool

    query := `
		SELECT EXISTS (
			SELECT 1 FROM ads WHERE uuid = $1
		)
	`

    err := r.db.QueryRowContext(ctx, query, uuid).Scan(&result)
    return result, err
}

func (repo *Repository) FindExpiredUuidList(ctx context.Context) ([]string, error) {
	var result = []string{}

	query := fmt.Sprintf(`
        SELECT
			uuid
		FROM 
			ads
		WHERE
			expires_at < now()
			and status = %d
    `, STATUS_ACTIVE)

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var uuid string
		if err := rows.Scan(&uuid); err != nil {
			return result, err
		}
		result = append(result, uuid)
	}

	return result, nil
}
