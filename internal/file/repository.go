package file

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type FileRepository struct {
	db *sql.DB
}

func NewFileRepository(db *sql.DB) *FileRepository {
	return &FileRepository{db}
}

func (r *FileRepository) Save(ctx context.Context, tx *sql.Tx, fileModel FileModel) error {
	query := `
        INSERT INTO files (
            ad_uuid,
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
		fileModel.AdUuid,
		fileModel.Path,
		fileModel.PreviewPath,
		fileModel.Order,
		fileModel.Size,
		fileModel.PreviewSize,
		fileModel.Mime,
		fileModel.PreviewMime,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *FileRepository) FindFilesByAdUuid(ctx context.Context, uuid uuid.UUID) ([]FileModel, error) {
	var result []FileModel

	query := `
		SELECT
			id,
			ad_uuid,
			path,
			preview_path,
			"order"
		FROM
			files
		WHERE 
			ad_uuid = $1
		ORDER BY
			"order" ASC
	`

	rows, err := r.db.QueryContext(ctx, query, uuid)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var file FileModel

		if err := rows.Scan(
			&file.Id,
			&file.AdUuid,
			&file.Path,
			&file.PreviewPath,
			&file.Order,
		); err != nil {
			return result, nil
		}

		result = append(result, file)
	}

	return result, nil
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
