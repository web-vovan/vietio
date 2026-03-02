package ads

import (
	"time"

	"github.com/google/uuid"
)

type AdsListQueryParams struct {
	Page       int
	CategoryId *int
	Sort       string
}

type AdsListFilterParams struct {
	Page       int
	Limit      int
	CategoryId *int
	Status     *int
	UserId     *int64
	Sort       string
	Order      string
}

type CreateAdRequestBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	CategoryId  int    `json:"category_id"`
}

type UpdateAdRequestBody struct {
	Uuid        uuid.UUID
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       int      `json:"price"`
	CategoryId  int      `json:"category_id"`
	OldImages   []string `json:"old_images"`
}

type AdModel struct {
	Uuid        uuid.UUID
	Title       string
	Description string
	UserId      int64
	CategoryId  int
	Price       int
	Status      int
	CreatedAt   time.Time
}

type AdsListRepository struct {
	Items []AdsListItemRepository
	Total int
}

type AdsListItemRepository struct {
	Uuid       uuid.UUID
	Title      string
	CategoryId int
	Price      int
	Status     int
	CreatedAt  time.Time
	Image      string
}

type AdsListItemResponse struct {
	Uuid       uuid.UUID `json:"uuid"`
	Title      string    `json:"title"`
	CategoryId int       `json:"category_id"`
	Price      int       `json:"price"`
	City       string    `json:"city"`
	Status     string    `json:"status"`
	Image      string    `json:"image"`
	CreatedAt  time.Time `json:"created_at"`
}

type AdsListResponse struct {
	Items []AdsListItemResponse `json:"items"`
	Total int                   `json:"total"`
	Limit int                   `json:"limit"`
	Page  int                   `json:"page"`
}

type MyAdsListResponse struct {
	Items []AdsListItemResponse `json:"items"`
	Total int                   `json:"total"`
}

type MySoldAdsListResponse struct {
	Items []AdsListItemResponse `json:"items"`
	Total int                   `json:"total"`
}

type MyFavoritesAdsListResponse struct {
	Items []AdsListItemResponse `json:"items"`
	Total int                   `json:"total"`
}

type CreateAdResponse struct {
	Uuid string `json:"uuid"`
}

type UpdateAdResponse struct {
	Result bool `json:"result"`
}

type DeleteAdResponse struct {
	Result bool `json:"result"`
}

type AdResponse struct {
	Uuid          uuid.UUID `json:"uuid"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	CategoryId    int       `json:"category_id"`
	Price         int       `json:"price"`
	City          string    `json:"city"`
	CreatedAt     time.Time `json:"created_at"`
	IsOwner       bool      `json:"is_owner"`
	IsFavorite    bool      `json:"is_favorite"`
	OwnerUsername string    `json:"owner_username"`
	Images        []string  `json:"images"`
}
