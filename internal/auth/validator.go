package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Validator struct{}

func NewValidator() *Validator {
    return &Validator{}
}

func (v *Validator) ValidateWebAppData(initData string, botToken string) (*TelegramUser, error) {
	// 1. Парсим строку запроса (query string)
	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse init data: %w", err)
	}

	// 2. Получаем хеш и удаляем его из списка параметров
	receivedHash := values.Get("hash")
	if receivedHash == "" {
		return nil, errors.New("hash is missing")
	}
	values.Del("hash")

	// 3. Проверка времени (защита от Replay Attack)
	// Данные валидны только в течение определенного времени (например, 24 часа)
	authDateStr := values.Get("auth_date")
	if authDateStr == "" {
		return nil, errors.New("auth_date is missing")
	}
	authDate, err := strconv.ParseInt(authDateStr, 10, 64)
	if err != nil {
		return nil, errors.New("invalid auth_date format")
	}

	if time.Now().Unix()-authDate > 86400 {
		return nil, errors.New("init data is expired")
	}

	// 4. Сортируем ключи по алфавиту
	var keys []string
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 5. Формируем строку проверки: key=value\nkey=value
	var dataCheckParts []string
	for _, k := range keys {
		dataCheckParts = append(dataCheckParts, fmt.Sprintf("%s=%s", k, values.Get(k)))
	}
	dataCheckString := strings.Join(dataCheckParts, "\n")

	// 6. Вычисляем секретный ключ
	// Secret key = HMAC_SHA256("WebAppData", BotToken)
	secretKey := hmacSha256([]byte("WebAppData"), []byte(botToken))

	// 7. Вычисляем хеш от строки проверки
	// Hash = HMAC_SHA256(SecretKey, DataCheckString)
	calculatedHash := hex.EncodeToString(hmacSha256(secretKey, []byte(dataCheckString)))

	// 8. Сравниваем полученный хеш с вычисленным
	if calculatedHash != receivedHash {
		return nil, errors.New("invalid hash signature")
	}

	// 9. Если всё ок — парсим пользователя
	userJSON := values.Get("user")
	var user TelegramUser
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		return nil, fmt.Errorf("failed to parse user json: %w", err)
	}

	return &user, nil
}

// Вспомогательная функция HMAC
func hmacSha256(key []byte, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}
