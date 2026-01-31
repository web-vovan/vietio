package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

// генерация тестовых данных initData
func generateTestInitData(botToken string, userName string) (string, error)  {
    type TelegramUser struct {
        ID              int64  `json:"id"`
        FirstName       string `json:"first_name"`
        LastName        string `json:"last_name"`
        Username        string `json:"username"`
        LanguageCode    string `json:"language_code"`
        AllowsWriteToPm bool   `json:"allows_write_to_pm"`
    }

    // Создаем фейкового пользователя
    user := TelegramUser{
        ID:              999999999,
        FirstName:       "Test",
        LastName:        "User",
        Username:        userName,
        LanguageCode:    "ru",
        AllowsWriteToPm: true,
    }

    // Сериализуем пользователя в JSON
    userJSON, err := json.Marshal(user)
    if err != nil {
        panic(err)
    }

    // Формируем список параметров (без hash)
    params := map[string]string{
        "query_id":  "AAGHsPI9AAAAAIew8j0swK3_",
        "user":      string(userJSON),
        "auth_date": fmt.Sprintf("%d", time.Now().Unix()),
    }

    // 3. Сортируем ключи по алфавиту для создания data_check_string
    keys := make([]string, 0, len(params))
    for k := range params {
        keys = append(keys, k)
    }
    sort.Strings(keys)

    // Собираем строку для проверки: key=value\nkey=value
    var dataCheckParts []string
    for _, k := range keys {
        dataCheckParts = append(dataCheckParts, fmt.Sprintf("%s=%s", k, params[k]))
    }
    dataCheckString := strings.Join(dataCheckParts, "\n")

    // Генерируем секретный ключ
    // Ключ: "WebAppData", Сообщение: BotToken
    secretKey := hmacSha256([]byte("WebAppData"), []byte(botToken))

    // Генерируем хеш
    // Ключ: secretKey, Сообщение: dataCheckString
    hash := hex.EncodeToString(hmacSha256(secretKey, []byte(dataCheckString)))

    // Формируем итоговую URL-строку (initData)
    // Используем url.Values для правильного кодирования (url encoding)
    values := url.Values{}
    for k, v := range params {
        values.Set(k, v)
    }
    values.Set("hash", hash)

    // Telegram передает параметры именно в таком виде (raw query string)
    finalInitData := values.Encode()

    return finalInitData , nil       
}

// Вспомогательная функция для HMAC-SHA256
func hmacSha256(key []byte, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}