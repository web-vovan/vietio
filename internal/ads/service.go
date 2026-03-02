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
	appErrors "vietio/internal/errors"
	fileApp "vietio/internal/file"
	"vietio/internal/user"

	"github.com/google/uuid"
)

var allowedSort = map[string]string{
	"date":  "created_at",
	"price": "price",
}

type Service struct {
	repo         *Repository
	fileRepo     FileRepository
	userRepo     UserRepository
	wishlistRepo WishlistRepository
	storage      FileStorage
	validator    *Validator
}

type FileRepository interface {
	Save(context.Context, *sql.Tx, fileApp.FileModel) error
	DeleteById(context.Context, *sql.Tx, int64) error
	FindFilesByAdUuid(context.Context, uuid.UUID) ([]fileApp.FileModel, error)
}

type UserRepository interface {
	GetUserById(context.Context, int64) (user.UserModel, error)
}

type WishlistRepository interface {
	AddWishlist(ctx context.Context, userId int64, adUuid uuid.UUID) error
	DeleteWishlist(ctx context.Context, userId int64, adUuid uuid.UUID) error
	HasUserWishlistByAdUuid(ctx context.Context, userId int64, adUuid uuid.UUID) (bool, error)
}

func NewService(
	repo *Repository,
	fileRepository FileRepository,
	userRepository UserRepository,
	wishlistRepository WishlistRepository,
	storage FileStorage,
	validator *Validator,
) *Service {
	return &Service{
		repo:         repo,
		fileRepo:     fileRepository,
		userRepo:     userRepository,
		wishlistRepo: wishlistRepository,
		storage:      storage,
		validator:    validator,
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
			Status:     getTextStatus(adItem.Status),
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

	filterParams := AdsListFilterParams{
		Page:   1,
		Sort:   "created_at",
		UserId: &userId,
		Order:  "desc",
		Limit:  1000,
	}
	adsListRepository, err := s.repo.FindAds(ctx, filterParams)
	if err != nil {
		return result, err
	}

	items := make([]AdsListItemResponse, 0, len(adsListRepository.Items))

	for _, adItem := range adsListRepository.Items {
		items = append(items, AdsListItemResponse{
			Uuid:       adItem.Uuid,
			Title:      adItem.Title,
			CategoryId: adItem.CategoryId,
			Price:      adItem.Price,
			City:       "Нячанг",
			Status:     getTextStatus(adItem.Status),
			Image:      s.storage.GetPublicPath(adItem.Image),
			CreatedAt:  adItem.CreatedAt,
		})
	}

	result.Items = items
	result.Total = len(items)

	return result, nil
}

func (s *Service) GetMySoldAds(ctx context.Context) (MySoldAdsListResponse, error) {
	var result MySoldAdsListResponse
	userId, err := authctx.GeUserIdFromContext(ctx)
	if err != nil {
		return result, err
	}

	status := STATUS_SOLD
	filterParams := AdsListFilterParams{
		Page:   1,
		Sort:   "created_at",
		UserId: &userId,
		Status: &status,
		Order:  "desc",
		Limit:  1000,
	}
	adsListRepository, err := s.repo.FindAds(ctx, filterParams)
	if err != nil {
		return result, err
	}

	items := make([]AdsListItemResponse, 0, len(adsListRepository.Items))

	for _, adItem := range adsListRepository.Items {
		items = append(items, AdsListItemResponse{
			Uuid:       adItem.Uuid,
			Title:      adItem.Title,
			CategoryId: adItem.CategoryId,
			Price:      adItem.Price,
			City:       "Нячанг",
			Status:     getTextStatus(adItem.Status),
			CreatedAt:  adItem.CreatedAt,
		})
	}

	result.Items = items
	result.Total = len(items)

	return result, nil
}

func (s *Service) GetMyFavoritesAds(ctx context.Context) (MyFavoritesAdsListResponse, error) {
	var result MyFavoritesAdsListResponse
	userId, err := authctx.GeUserIdFromContext(ctx)
	if err != nil {
		return result, err
	}

	adsListRepository, err := s.repo.FindFavoritesAdsByUserId(ctx, userId)
	if err != nil {
		return result, err
	}

	items := make([]AdsListItemResponse, 0, len(adsListRepository.Items))

	for _, adItem := range adsListRepository.Items {
		status := getTextStatus(adItem.Status)
		var image string

		if status == "active" {
			image = s.storage.GetPublicPath(adItem.Image)
		}

		items = append(items, AdsListItemResponse{
			Uuid:       adItem.Uuid,
			Title:      adItem.Title,
			CategoryId: adItem.CategoryId,
			Price:      adItem.Price,
			City:       "Нячанг",
			Status:     status,
			Image:      image,
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

	if adModel.Status != STATUS_ACTIVE {
		return result, appErrors.ErrAdNotActive
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

	isFavorite, err := s.wishlistRepo.HasUserWishlistByAdUuid(ctx, ctxUserId, uuid)
	if err != nil {
		return result, appErrors.ErrAdFavorite
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
		IsFavorite:    isFavorite,
		OwnerUsername: adOwner.Username,
		Images:        images,
	}, nil
}

func (s *Service) UpdateAd(ctx context.Context, payload UpdateAdRequestBody, images []*multipart.FileHeader) (UpdateAdResponse, error) {
	result := UpdateAdResponse{}

	contextUserId, err := authctx.GeUserIdFromContext(ctx)
	if err != nil {
		return result, err
	}

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

	if ad.UserId != contextUserId {
		return result, appErrors.ErrForbidden
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

func (s *Service) ArchiveAd(ctx context.Context, uuid uuid.UUID) error {
	ad, err := s.repo.FindAdByUuid(ctx, uuid)
	if err != nil {
		return err
	}

	return s.processDeleteAd(ctx, ad, STATUS_EXPIRED)
}

func (s *Service) DeleteAd(ctx context.Context, uuid uuid.UUID) error {
	contextUserId, err := authctx.GeUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	ad, err := s.repo.FindAdByUuid(ctx, uuid)
	if err != nil {
		return err
	}

	if ad.UserId != contextUserId {
		return appErrors.ErrForbidden
	}

	return s.processDeleteAd(ctx, ad, STATUS_USER_DELETED)
}

func (s *Service) processDeleteAd(ctx context.Context, ad AdModel, finalStatus int) error {
	tx, err := s.repo.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	files, err := s.fileRepo.FindFilesByAdUuid(ctx, ad.Uuid)
	if err != nil {
		return err
	}

	for _, f := range files {
		err = s.fileRepo.DeleteById(ctx, tx, f.Id)
		if err != nil {
			return err
		}

		err = s.storage.DeleteByPath(ctx, f.Path)
		if err != nil {
			return err
		}

		err = s.storage.DeleteByPath(ctx, f.PreviewPath)
		if err != nil {
			return err
		}
	}

	err = s.repo.ChangeStatusAdByUuidWithTx(ctx, tx, finalStatus, ad.Uuid)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return nil
	}

	return err
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

func (s *Service) ArchivingAds(ctx context.Context) error {
	uuidList, err := s.repo.FindExpiredUuidList(ctx)
	if err != nil {
		return err
	}

	for _, uuidItem := range uuidList {
		r, err := uuid.Parse(uuidItem)
		if err != nil {
			return err
		}
		s.ArchiveAd(ctx, r)
	}

	return nil
}

func (s *Service) MarkingSoldAd(ctx context.Context, uuid uuid.UUID) error {
	contextUserId, err := authctx.GeUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	ad, err := s.repo.FindAdByUuid(ctx, uuid)
	if err != nil {
		return err
	}

	if ad.UserId != contextUserId {
		return appErrors.ErrForbidden
	}

	return s.processDeleteAd(ctx, ad, STATUS_SOLD)
}

func (s *Service) AddFavorite(ctx context.Context, uuid uuid.UUID) error {
	contextUserId, err := authctx.GeUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	err = s.wishlistRepo.AddWishlist(ctx, contextUserId, uuid)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteFavorite(ctx context.Context, uuid uuid.UUID) error {
	contextUserId, err := authctx.GeUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	err = s.wishlistRepo.DeleteWishlist(ctx, contextUserId, uuid)
	if err != nil {
		return err
	}

	return nil
}
