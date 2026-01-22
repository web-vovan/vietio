package ads

import (
	"context"
	"database/sql"
	"fmt"
	"mime/multipart"
	"strings"
	filePkg "vietio/internal/file"
)

var allowedSort = map[string]string{
	"date":  "created_at",
	"price": "price",
}

type Service struct {
	repo            *Repository
	categoryChecker CategoryChecker
	storage         FileStorage
	fileRepository  FileRepository
}

type CategoryChecker interface {
	Exists(context.Context, int) (bool, error)
}

type FileRepository interface {
	Save(context.Context, *sql.Tx, filePkg.File) error
}

func NewService(
	repo *Repository,
	categoryChecker CategoryChecker,
	storage FileStorage,
	fileRepository FileRepository,
) *Service {
	return &Service{
		repo:            repo,
		categoryChecker: categoryChecker,
		storage:         storage,
		fileRepository:  fileRepository,
	}
}

func (s *Service) GetAds(ctx context.Context, params AdsListQueryParams) (AdsListResponse, error) {
	var categoryId *int
	page := 1
	sort := "created_at"
	order := "desc"

	if params.Page > 0 {
		page = params.Page
	}

	sortParts := strings.Split(params.Sort, "_")

	if len(sortParts) == 2 {
		if v, ok := allowedSort[sortParts[0]]; ok {
			sort = v
		}
		if sortParts[1] == "asc" {
			order = "asc"
		}
	}

	if params.CategoryId != nil && *params.CategoryId > 0 {
		categoryId = params.CategoryId
	}

	filterParams := AdsListFilterParams{
		Page:       page,
		CategoryId: categoryId,
		Sort:       sort,
		Order:      order,
		Limit:      20,
	}

	adsListDB, _ := s.repo.FindAds(ctx, filterParams)

	return AdsListResponse{
		Items: adsListDB.Items,
		Total: adsListDB.Total,
		Limit: filterParams.Limit,
		Page:  filterParams.Page,
	}, nil
}

func (s *Service) CreateAd(ctx context.Context, payload CreateAdRequestBody, files []*multipart.FileHeader) (CreateAdResponse, error) {
	result := CreateAdResponse{}

	errors, err := s.validate(ctx, payload, files)
	if err != nil {
		return result, err
	}
	if len(errors) > 0 {
		return result, fmt.Errorf(strings.Join(errors, ";"))
	}

	tx, err := s.repo.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return result, err
	}
	defer tx.Rollback()

	id, err := s.repo.CreateAd(ctx, tx, payload)
	if err != nil {
		return result, fmt.Errorf("возникла ошибка при сохранении объявления: %w", err)
	}

	for i, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return result, err
		}
		defer file.Close()

		fileInfo, err := s.storage.Save(ctx, file, fileHeader)
		if err != nil {
			return result, err
		}

		fileModel := filePkg.File{
			AdId:  id,
			Path:  fileInfo.FileName,
			Order: i + 1,
			Mime:  fileInfo.Mime,
			Size:  fileInfo.Size,
		}

		err = s.fileRepository.Save(ctx, tx, fileModel)
		if err != nil {
			return result, err
		}
	}

	result.Id = id

	if err := tx.Commit(); err != nil {
		return result, nil
	}

	return result, err
}

func (s *Service) validate(
	ctx context.Context,
	payload CreateAdRequestBody,
	files []*multipart.FileHeader,
) ([]string, error) {
	var errors []string

	if payload.Title == "" {
		errors = append(errors, "title не может быть пустым")
	}
	if payload.Description == "" {
		errors = append(errors, "description не может быть пустым")
	}
	if payload.Price < 0 {
		errors = append(errors, "price не может быть отрицательным")
	}
	if payload.CategoryId < 1 {
		errors = append(errors, "category_id должен быть >= 1")
	}
	if payload.CategoryId >= 1 {
		exists, err := s.categoryChecker.Exists(ctx, payload.CategoryId)
		if err != nil {
			return nil, fmt.Errorf("ошибка БД при проверки существования категории")
		}
		if !exists {
			errors = append(errors, "category_id такой категории не существует")
		}
	}
	if len(files) == 0 || len(files) > 3 {
		errors = append(errors, "files должен быть > 0 и меньше 3")
	}

	return errors, nil
}
