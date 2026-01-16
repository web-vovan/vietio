package ads

import (
	"context"
	"strings"
)

var allowedSort = map[string]string{
    "date": "created_at",
    "price": "price",
}

type Service struct {
    repo *Repository
}

func NewService(repo *Repository) *Service {
    return &Service{
        repo: repo,
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
        Page: page,
        CategoryId: categoryId,
        Sort: sort,
        Order: order,
        Limit: 20,
    }

    adsListDB, _ := s.repo.FindAds(ctx, filterParams)

    return AdsListResponse{
        Items: adsListDB.Items,
        Total: adsListDB.Total,
        Limit: filterParams.Limit,
        Page: filterParams.Page,
    }, nil
}
