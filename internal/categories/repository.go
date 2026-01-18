package categories

import (
	"context"
	"database/sql"
)

type Repository struct {
    db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
    return &Repository{
        db: db,
    }
}

func (r *Repository) Exists(ctx context.Context, categoryId int) (bool, error) {
    var result bool

    query := `
		SELECT EXISTS (
			SELECT 1 FROM categories WHERE id = $1
		)
	`

    err := r.db.QueryRowContext(ctx, query, categoryId).Scan(&result)
    return result, err
}