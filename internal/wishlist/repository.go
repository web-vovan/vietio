package wishlist

import (
	"context"
	"database/sql"

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

func (r *Repository) AddWishlist(ctx context.Context, userId int64, adUuid uuid.UUID) error {
    query := `
        INSERT INTO wishlist (user_id, ad_uuid)
        VALUES ($1, $2)
    `

    _, err := r.db.ExecContext(ctx, query, userId, adUuid)
    if err != nil {
        return err
    }

    return nil
}

func (r *Repository) DeleteWishlist(ctx context.Context, userId int64, adUuid uuid.UUID) error {
    query := `
        DELETE FROM wishlist
        WHERE user_id=$1 and ad_uuid=$2
    `

    _, err := r.db.ExecContext(ctx, query, userId, adUuid)
    if err != nil {
        return err
    }

    return nil
}