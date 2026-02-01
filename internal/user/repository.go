package user

import (
	"context"
	"database/sql"
	"errors"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetUserByTelegramId(ctx context.Context, telegramId int64) (UserModel, error) {
    var user UserModel

    query := `
        SELECT
            id,
            telegram_id,
            username
        FROM
            users
        WHERE
            telegram_id = $1
        LIMIT 1
    `

    err := r.db.QueryRowContext(ctx, query, telegramId).Scan(
        &user.Id,
        &user.TelegramId,
        &user.Username,
    )

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return user, sql.ErrNoRows
        }
        return user, err
    }

    return user, nil
}

func (r *Repository) UpdateUsername(ctx context.Context, user UserModel) error {
    query := `
        UPDATE 
            users
        SET 
            username = $1,
            updated_at = now()
        WHERE 
            telegram_id = $2
    `

    _, err := r.db.ExecContext(ctx, query, user.Username, user.TelegramId)
    if err != nil {
        return err
    }

    return nil
}

func (r *Repository) CreateUser(ctx context.Context, user UserModel) (int64, error) {
    var id int64

	query := `
		INSERT INTO users (
			telegram_id,
			username
		)
		VALUES ($1, $2)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.TelegramId,
		user.Username,
	).Scan(&id)

	if err != nil {
		return id, err
	}

	return id, nil
}