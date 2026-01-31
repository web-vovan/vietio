package auth

import (
	"encoding/json"
	"net/http"

	"vietio/internal/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetToken(w http.ResponseWriter, r *http.Request) {
	payload := AuthLoginRequestBody{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := h.service.GetJwtToken(r.Context(), payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Json(w, result, http.StatusOK)
}

func (h *Handler) GetTestInitData(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	if username == "" {
		http.Error(w, "username не может быть пустым", http.StatusInternalServerError)
		return
	}

	result, err := h.service.GenerateTestInitData(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Json(w, result, http.StatusOK)
}

