package auth

import "github.com/golang-jwt/jwt/v5"

type TelegramUser struct {
	ID              int64  `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Username        string `json:"username"`
	LanguageCode    string `json:"language_code"`
	AllowsWriteToPm bool   `json:"allows_write_to_pm"`
}

type AuthLoginRequestBody struct {
	InitData string `json:"init_data"`
}

type AuthLoginResponse struct {
	Token string `json:"token"`
}

type TestInitDataResponse struct {
	InitData string `json:"init_data"`
}

type AccessTokenClaims struct {
	UserId           int64 `json:"user_id"`
	TelegramId       int64 `json:"telegram_id"`
	jwt.RegisteredClaims
}
