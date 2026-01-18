package ads

import (
	"context"
	"fmt"
	"strings"
)

var allowedSort = map[string]string{
	"date":  "created_at",
	"price": "price",
}

type Service struct {
	repo            *Repository
	categoryChecker CategoryChecker
}

type CategoryChecker interface {
	Exists(context.Context, int) (bool, error)
}

func NewService(repo *Repository, categoryChecker CategoryChecker) *Service {
	return &Service{
		repo:            repo,
		categoryChecker: categoryChecker,
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

func (s *Service) CreateAd(ctx context.Context, payload CreateAdRequestBody) (CreateAdResponse, error) {
	result := CreateAdResponse{}
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
			return result, fmt.Errorf("ошибка БД при проверки существования категории")
		}
		if !exists {
            errors = append(errors, "category_id такой категории не существует")
		}
	}

	if len(errors) > 0 {
		return result, fmt.Errorf(strings.Join(errors, ";"))
	}

    id, err := s.repo.CreateAd(ctx, payload)
    if err != nil {
        return result, fmt.Errorf("возникла ошибка при сохранении объявления: %w", err)
    }
	result.Id = id

	return result, nil
}
