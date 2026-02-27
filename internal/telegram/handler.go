package telegram

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"vietio/internal/response"
)

type Handler struct {
	Logger   *slog.Logger
	TgClient *Client
}

func NewHandler(logger *slog.Logger, tgClient *Client) *Handler {
	return &Handler{
		Logger:   logger,
		TgClient: tgClient,
	}
}

func (h *Handler) Webhook(w http.ResponseWriter, r *http.Request) {
	payload := WebhookRequestBody{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		h.Logger.Warn("telegram webhook error", "err", err)
		response.Json(w, "", http.StatusOK)
		return
	}

	if payload.Message.Text != "/start" {
		h.Logger.Info("telegram webhook not start", "request", r.Body)
		response.Json(w, "", http.StatusOK)
		return
	}

	go func(chatId int64) {
        text := `Добро пожаловать в Vietio 🇻🇳

Локальные объявления в Нячанге — всё в одном месте.

👇 Нажмите «Открыть» и начните пользоваться`

		err := h.TgClient.SendMessage(chatId, text)
        if err != nil {
            h.Logger.Warn("sendMessage telegram error", "err", err)
        }
	}(payload.Message.Chat.Id)

	response.Json(w, "", http.StatusOK)
}
