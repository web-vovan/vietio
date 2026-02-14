package ads

import (
	"errors"
	"log/slog"
	"net/http"

	appErrors "vietio/internal/errors"
	"vietio/internal/response"
	"vietio/pkg/utils"

	"github.com/google/uuid"
)

type Handler struct {
	service *Service
	logger  *slog.Logger
}

func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
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
		h.logger.Error(appErrors.ErrAdsList.Error(), "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response.Json(w, result, http.StatusOK)
}

func (h *Handler) GetAd(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		http.Error(w, "невалидный uuid в запросе", http.StatusBadRequest)
		return
	}

	result, err := h.service.GetAd(r.Context(), uuid)
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrAdNotFound):
			http.Error(w, appErrors.ErrAdNotFound.Error(), http.StatusNotFound)
		default:
			h.logger.Error(appErrors.ErrAd.Error(), "err", err, "uuid", uuid)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	response.Json(w, result, http.StatusOK)
}

func (h *Handler) GetMyAds(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.GetMyAds(r.Context())
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

	validationErrors := appErrors.NewValidationError()

	price := validateIntField("price", r.FormValue("price"), false, 0, validationErrors)
	categoryId := validateIntField("category_id", r.FormValue("category_id"), true, 0, validationErrors)

	if validationErrors.HasErrors() {
		h.logger.Warn(appErrors.ErrCreateAdValidation.Error(), "err", err)
		response.Json(w, validationErrors, http.StatusBadRequest)
		return
	}

	payload := CreateAdRequestBody{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Price:       price,
		CategoryId:  categoryId,
	}

	images := r.MultipartForm.File["images"]

	result, err := h.service.CreateAd(r.Context(), payload, images)

	if err != nil {
		var vError *appErrors.ValidationError
		if errors.As(err, &vError) {
			h.logger.Warn(appErrors.ErrCreateAdValidation.Error(), "err", err, "payload", payload)
			response.Json(w, err, http.StatusBadRequest)
		} else {
			h.logger.Error(appErrors.ErrCreateAd.Error(), "err", err, "payload", payload)
			http.Error(w, "internal server", http.StatusInternalServerError)
		}
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

	validationErrors := appErrors.NewValidationError()

	price := validateIntField("price", r.FormValue("price"), false, 0, validationErrors)
	categoryId := validateIntField("category_id", r.FormValue("category_id"), true, 0, validationErrors)

	if validationErrors.HasErrors() {
		h.logger.Warn(appErrors.ErrCreateAdValidation.Error(), "err", err)
		response.Json(w, validationErrors, http.StatusBadRequest)
		return
	}

	payload := UpdateAdRequestBody{
		Uuid:        uuid,
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Price:       price,
		CategoryId:  categoryId,
		OldImages:   r.Form["old_images"],
	}

	images := r.MultipartForm.File["images"]

	result, err := h.service.UpdateAd(r.Context(), payload, images)

	if err != nil {
		var vError *appErrors.ValidationError
		if errors.As(err, &vError) {
			h.logger.Warn(appErrors.ErrUpdateAdValidation.Error(), "err", err, "payload", payload)
			response.Json(w, err, http.StatusBadRequest)
		} else if errors.Is(err, appErrors.ErrForbidden) {
			h.logger.Warn(appErrors.ErrForbidden.Error(), "err", err, "payload", payload)
			http.Error(w, "forbidden", http.StatusForbidden)
		} else {
			h.logger.Error(appErrors.ErrUpdateAd.Error(), "err", err, "payload", payload)
			http.Error(w, "internal server", http.StatusInternalServerError)
		}
		return
	}

	result.Result = true

	response.Json(w, result, http.StatusOK)
}

func (h *Handler) DeleteAd(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		http.Error(w, "невалидный uuid в запросе", http.StatusInternalServerError)
		return
	}

	result, err := h.service.DeleteAd(r.Context(), uuid)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Json(w, result, http.StatusOK)
}
