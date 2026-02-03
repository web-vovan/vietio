package ads

import (
	"errors"
	"net/http"

	"vietio/internal/authctx"
	appErrors "vietio/internal/errors"
	"vietio/internal/response"
	"vietio/pkg/utils"

	"github.com/google/uuid"
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
	// Максимальный размер тела 20 MB
	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		response.Json(w, err.Error(), http.StatusBadRequest)
		return
	}

	payload := CreateAdRequestBody{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Price:       utils.ParseInt(r.FormValue("price"), 0),
		CategoryId:  utils.ParseInt(r.FormValue("category_id"), 0),
	}

	images := r.MultipartForm.File["images"]

	result, err := h.service.CreateAd(r.Context(), payload, images)

	if err != nil {
		var validationError *appErrors.ValidationError
		if errors.As(err, &validationError) {
			response.Json(w, err, http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response.Json(w, result, http.StatusOK)
}

func (h *Handler) GetAd(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		http.Error(w, "невалидный uuid в запросе", http.StatusInternalServerError)
		return
	}

	result, err := h.service.GetAd(r.Context(), uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Json(w, result, http.StatusOK)
}

func (h *Handler) UpdateAd(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		http.Error(w, "невалидный uuid в запросе", http.StatusInternalServerError)
		return
	}

	// Максимальный размер тела 20 MB
	err = r.ParseMultipartForm(20 << 20)
	if err != nil {
		response.Json(w, err.Error(), http.StatusBadRequest)
		return
	}

	payload := UpdateAdRequestBody{
		Uuid:        uuid,
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Price:       utils.ParseInt(r.FormValue("price"), 0),
		CategoryId:  utils.ParseInt(r.FormValue("category_id"), 0),
		OldImages:   r.Form["old_images"],
	}

	images := r.MultipartForm.File["images"]

	result, err := h.service.UpdateAd(r.Context(), payload, images)

	if err != nil {
		var validationError *appErrors.ValidationError
		if errors.As(err, &validationError) {
			response.Json(w, err, http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response.Json(w, result, http.StatusOK)
}
