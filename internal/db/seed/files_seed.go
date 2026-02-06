package seed

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/textproto"

	"path/filepath"
	"vietio/internal/ads"
)

type nopMultipartFile struct {
	*bytes.Reader
}

func (f nopMultipartFile) Close() error {
	return nil
}

func runFilesSeed(dbConn *sql.DB, fileStorage ads.FileStorage, fsys fs.FS) error {
	ctx := context.Background()
	entries, err := fs.ReadDir(fsys, "images")
	if err != nil {
		return err
	}

	fileList := make([]*ads.FileInfo, 0, len(entries))

	// заливаем файлы в хранилище
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		filePath := filepath.Join("images", entry.Name())

		data, err := fs.ReadFile(fsys, filePath)
		if err != nil {
			return err
		}

		file := nopMultipartFile{
			Reader: bytes.NewReader(data),
		}

		mimeType := http.DetectContentType(data)

		header := &multipart.FileHeader{
			Filename: fileName,
			Size:     int64(len(data)),
			Header: textproto.MIMEHeader{
				"Content-Type": []string{mimeType},
			},
		}
		uploadedFile, err := fileStorage.Save(ctx, file, header)
		if err != nil {
			return err
		}

		fileList = append(fileList, uploadedFile)
	}

	adsUuids, err := getAllUuidFromTable(dbConn, "ads")
	if err != nil {
		return err
	}

	query := `
		INSERT INTO files (
			ad_uuid,
			path, 
			preview_path, 
			size, 
			preview_size, 
			mime, 
			preview_mime,
			storage
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	stmt, err := dbConn.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, uuid := range adsUuids {
		randInd := rand.Intn(len(fileList))

		_, err := stmt.Exec(
			uuid,
			fileList[randInd].FileName,
			fileList[randInd].PreviewFileName,
			fileList[randInd].Size,
			fileList[randInd].PreviewSize,
			fileList[randInd].Mime,
			fileList[randInd].PreviewMime,
			"s3",
		)
		if err != nil {
			return err
		}
	}

	return nil
}
