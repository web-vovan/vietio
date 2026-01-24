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

func (r *FileRepository) Save(ctx context.Context, tx *sql.Tx, file File) error {
	query := `
        INSERT INTO files (
            ad_id,
            path,
            preview_path,
            "order",
            size,
            preview_size,
            mime,
            preview_mime
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8
        )
    `

	_, err := tx.ExecContext(ctx, query,
		file.AdId,
		file.Path,
		file.PreviewPath,
		file.Order,
		file.Size,
		file.PreviewSize,
		file.Mime,
		file.PreviewMime,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *FileRepository) DeleteByPath(ctx context.Context, path string) error {
	query := `
        DELETE FROM files 
        WHERE path = $1
    `

	_, err := r.db.Exec(query, path)
	if err != nil {
		return err
	}

	return nil
}
