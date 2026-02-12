package ads

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mime/multipart"
	"net/url"
	"path"
	"strings"

	"vietio/internal/authctx"
	fileApp "vietio/internal/file"
	"vietio/internal/user"
	appErrors "vietio/internal/errors"

	"github.com/google/uuid"
)

var allowedSort = map[string]string{
	"date":  "created_at",
	"price": "price",
}

type Service struct {
	repo      *Repository
	fileRepo  FileRepository
	userRepo  UserRepository
	storage   FileStorage
	validator *Validator
}

type FileRepository interface {
	Save(context.Context, *sql.Tx, fileApp.FileModel) error
	DeleteById(context.Context, *sql.Tx, int64) error
	FindFilesByAdUuid(context.Context, uuid.UUID) ([]fileApp.FileModel, error)
}

type UserRepository interface {
	GetUserById(context.Context, int64) (user.UserModel, error)
}

func NewService(
	repo *Repository,
	fileRepository FileRepository,
	userRepository UserRepository,
	storage FileStorage,
	validator *Validator,
) *Service {
	return &Service{
		repo:      repo,
		fileRepo:  fileRepository,
		userRepo:  userRepository,
		storage:   storage,
		validator: validator,
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

	adsListRepository, err := s.repo.FindAds(ctx, filterParams)
	if err != nil {
		return AdsListResponse{}, err
	}

	items := make([]AdsListItemResponse, 0, len(adsListRepository.Items))

	for _, adItem := range adsListRepository.Items {
		items = append(items, AdsListItemResponse{
			Uuid:       adItem.Uuid,
			Title:      adItem.Title,
			CategoryId: adItem.CategoryId,
			Price:      adItem.Price,
			City:       "Нячанг",
			Image:      s.storage.GetPublicPath(adItem.Image),
			CreatedAt:  adItem.CreatedAt,
		})
	}

	return AdsListResponse{
		Items: items,
		Total: adsListRepository.Total,
		Limit: filterParams.Limit,
		Page:  filterParams.Page,
	}, nil
}

func (s *Service) GetMyAds(ctx context.Context) (MyAdsListResponse, error) {
	var result MyAdsListResponse
	userId, err := authctx.GeUserIdFromContext(ctx)
	if err != nil {
		return result, err
	}

	if userId == 0 {
		return result, errors.New("пользователь не авторизован")
	}

	filterParams := AdsListFilterParams{
		Page:   1,
		Sort:   "created_at",
		UserId: &userId,
		Order:  "desc",
		Limit:  1000,
	}
	adsListRepository, _ := s.repo.FindAds(ctx, filterParams)

	items := make([]AdsListItemResponse, 0, len(adsListRepository.Items))

	for _, adItem := range adsListRepository.Items {
		items = append(items, AdsListItemResponse{
			Uuid:       adItem.Uuid,
			Title:      adItem.Title,
			CategoryId: adItem.CategoryId,
			Price:      adItem.Price,
			City:       "Нячанг",
			Image:      s.storage.GetPublicPath(adItem.Image),
			CreatedAt:  adItem.CreatedAt,
		})
	}

	result.Items = items
	result.Total = len(items)

	return result, nil
}

func (s *Service) CreateAd(ctx context.Context, payload CreateAdRequestBody, images []*multipart.FileHeader) (CreateAdResponse, error) {
	result := CreateAdResponse{}

	validationErrors := s.validator.createAdValidate(ctx, payload, images)
	if validationErrors.HasErrors() {
		return result, validationErrors
	}

	tx, err := s.repo.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return result, err
	}
	defer tx.Rollback()

	uuid, err := s.repo.CreateAd(ctx, tx, payload)
	if err != nil {
		return result, fmt.Errorf("возникла ошибка при сохранении объявления: %w", err)
	}

	err = s.saveNewImages(ctx, tx, uuid, images)
	if err != nil {
		return result, err
	}

	result.Uuid = uuid.String()

	if err := tx.Commit(); err != nil {
		return result, nil
	}

	return result, err
}

func (s *Service) GetAd(ctx context.Context, uuid uuid.UUID) (AdResponse, error) {
	var result AdResponse

	ctxUserId, err := authctx.GeUserIdFromContext(ctx)
	if err != nil {
		return result, err
	}

	adModel, err := s.repo.FindAdByUuid(ctx, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, appErrors.ErrAdNotFound
		}
		return result, err
	}

	adFiles, err := s.fileRepo.FindFilesByAdUuid(ctx, uuid)
	if err != nil {
		return result, err
	}

	adOwner, err := s.userRepo.GetUserById(ctx, adModel.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, appErrors.ErrAdUserNotFound
		}
		return result, err
	}

	var images = make([]string, 0, len(adFiles))
	for _, file := range adFiles {
		publicPath := s.storage.GetPublicPath(file.Path)
		images = append(images, publicPath)
	}

	return AdResponse{
		Uuid:          adModel.Uuid,
		Title:         adModel.Title,
		Description:   adModel.Description,
		CategoryId:    adModel.CategoryId,
		Price:         adModel.Price,
		City:          "Нячанг",
		CreatedAt:     adModel.CreatedAt,
		IsOwner:       adModel.UserId == ctxUserId,
		OwnerUsername: adOwner.Username,
		Images:        images,
	}, nil
}

func (s *Service) UpdateAd(ctx context.Context, payload UpdateAdRequestBody, images []*multipart.FileHeader) (UpdateAdResponse, error) {
	result := UpdateAdResponse{}

	validationErrors := s.validator.updateAdValidate(ctx, payload, images)
	if validationErrors.HasErrors() {
		return result, validationErrors
	}

	tx, err := s.repo.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return result, err
	}
	defer tx.Rollback()

	ad, err := s.repo.FindAdByUuid(ctx, payload.Uuid)
	if err != nil {
		return result, err
	}

	ad.Title = payload.Title
	ad.Description = payload.Description
	ad.Price = payload.Price
	ad.CategoryId = payload.CategoryId

	err = s.repo.UpdateAd(ctx, tx, ad)
	if err != nil {
		return result, err
	}

	var oldImagesMap = make(map[string]bool)

	for _, i := range payload.OldImages {
		u, err := url.Parse(i)
		if err != nil {
			return result, err
		}
		oldImagesMap[path.Base(u.Path)] = true
	}

	oldFiles, err := s.fileRepo.FindFilesByAdUuid(ctx, payload.Uuid)
	if err != nil {
		return result, err
	}

	var filesToDelete []fileApp.FileModel

	for _, f := range oldFiles {
		if _, ok := oldImagesMap[f.Path]; !ok {
			filesToDelete = append(filesToDelete, f)
		}
	}

	for _, f := range filesToDelete {
		err = s.fileRepo.DeleteById(ctx, tx, f.Id)
		if err != nil {
			return result, err
		}

		err = s.storage.DeleteByPath(ctx, f.Path)
		if err != nil {
			return result, err
		}

		err = s.storage.DeleteByPath(ctx, f.PreviewPath)
		if err != nil {
			return result, err
		}
	}

	err = s.saveNewImages(ctx, tx, payload.Uuid, images)
	if err != nil {
		return result, err
	}

	if err := tx.Commit(); err != nil {
		return result, nil
	}

	return result, err
}

func (s *Service) DeleteAd(ctx context.Context, uuid uuid.UUID) (DeleteAdResponse, error) {
	var result DeleteAdResponse

	contextUserId, err := authctx.GeUserIdFromContext(ctx)
	if err != nil {
		return result, err
	}

	tx, err := s.repo.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return result, err
	}
	defer tx.Rollback()

	ad, err := s.repo.FindAdByUuid(ctx, uuid)
	if err != nil {
		return result, err
	}

	if ad.UserId != contextUserId {
		return result, errors.New("нет прав")
	}

	files, err := s.fileRepo.FindFilesByAdUuid(ctx, uuid)
	if err != nil {
		return result, err
	}

	for _, f := range files {
		err = s.fileRepo.DeleteById(ctx, tx, f.Id)
		if err != nil {
			return result, err
		}

		err = s.storage.DeleteByPath(ctx, f.Path)
		if err != nil {
			return result, err
		}

		err = s.storage.DeleteByPath(ctx, f.PreviewPath)
		if err != nil {
			return result, err
		}
	}

	err = s.repo.DeleteAdByUuid(ctx, tx, uuid)
	if err != nil {
		return result, err
	}

	result.Result = true

	return result, err
}

func (s *Service) saveNewImages(
	ctx context.Context,
	tx *sql.Tx,
	adUuid uuid.UUID,
	images []*multipart.FileHeader,
) error {
	for _, fileHeader := range images {
		file, err := fileHeader.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		fileInfo, err := s.storage.Save(ctx, file, fileHeader)
		if err != nil {
			return err
		}

		fileModel := fileApp.FileModel{
			AdUuid:      adUuid,
			Path:        fileInfo.FileName,
			PreviewPath: fileInfo.PreviewFileName,
			Mime:        fileInfo.Mime,
			PreviewMime: fileInfo.PreviewMime,
			Size:        fileInfo.Size,
			PreviewSize: fileInfo.PreviewSize,
			Storage:     s.storage.GetType(),
		}

		err = s.fileRepo.Save(ctx, tx, fileModel)
		if err != nil {
			return err
		}
	}

	return nil
}
