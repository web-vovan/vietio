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

type Ad struct {
	Uuid        uuid.UUID
	Title       string
	Description string
	CategoryId  int
	Price       int
	CreatedAt   time.Time
}

type AdsListDB struct {
	Items []Ad
	Total int
}

type AdsListResponse struct {
	Items []Ad `json:"items"`
	Total int  `json:"total"`
	Limit int  `json:"limit"`
	Page  int  `json:"page"`
}

type CreateAdRequestBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	CategoryId  int    `json:"category_id"`
}

type CreateAdResponse struct {
	Uuid string `json:"uuid"`
}
