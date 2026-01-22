package file

import (
	"context"
	"database/sql"
)

type FileRepository struct {
	db *sql.DB
}

func NewFileRepository(db *sql.DB) *FileRepository {
	return &FileRepository{db}
}

func (r *FileRepository) Save(ctx context.Context, tx *sql.Tx,  file File) error {
    query := `
        INSERT INTO files (
            ad_id,
            path,
            "order",
            size,
            mime
        ) VALUES (
            $1, $2, $3, $4, $5
        )
    `

    _, err := tx.ExecContext(ctx, query, file.AdId, file.Path, file.Order, file.Size, file.Mime)
    if err != nil {
        return err
    }

	return nil
}
