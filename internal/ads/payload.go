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
	Sort       string
	Order      string
}

type CreateAdRequestBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	CategoryId  int    `json:"category_id"`
}

type AdModel struct {
	Uuid        uuid.UUID
	Title       string
	Description string
	CategoryId  int
	Price       int
	CreatedAt   time.Time
}

type AdsListRepository struct {
	Items []AdModel
	Total int
}

type AdsListItemResponse struct {
	Uuid       uuid.UUID `json:"uuid"`
	Title      string    `json:"title"`
	CategoryId int       `json:"category_id"`
	Price      int       `json:"price"`
	City       string    `json:"city"`
	Image      string    `json:"image"`
	CreatedAt  time.Time `json:"created_at"`
}

type AdsListResponse struct {
	Items []AdsListItemResponse `json:"items"`
	Total int                   `json:"total"`
	Limit int                   `json:"limit"`
	Page  int                   `json:"page"`
}

type CreateAdResponse struct {
	Uuid string `json:"uuid"`
}

type AdResponse struct {
	Uuid        uuid.UUID `json:"uuid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CategoryId  int       `json:"category_id"`
	Price       int       `json:"price"`
	City        string    `json:"city"`
	CreatedAt   time.Time `json:"created_at"`
	Images      []string  `json:"images"`
}
