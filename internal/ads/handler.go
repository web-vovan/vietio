package ads

import (
	"encoding/json"
	"net/http"
	"vietio/internal/response"
	"vietio/pkg/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetAds(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	params := AdsListQueryParams{
		Page:       utils.ParseInt(q.Get("page"), 1),
		CategoryId: utils.ParseNullableInt(q.Get("category_id")),
		Sort:       q.Get("sort"),
	}

	result, err := h.service.GetAds(r.Context(), params)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
		return
    }

	response.Json(w, result, http.StatusOK)
}

func (h *Handler) CreateAd(w http.ResponseWriter, r *http.Request) {
	payload := CreateAdRequestBody{}


	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		response.Json(w, err.Error(), http.StatusBadRequest)
		return
    }

	type Test struct {
		Id int `json:"id"`
	}

	t := Test{
		Id: 123,
	}


	response.Json(w, t, http.StatusOK)
}
