package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	Token  string
	Client *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		Token: token,
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func(c *Client) SendMessage(chatId int64, msg string) error {
    url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.Token)

    payload := map[string]any{
        "chat_id": chatId,
        "text": msg,
    }

    body, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    resp, err := c.Client.Post(
        url,
        "application/json",
        bytes.NewBuffer(body),
    )
    if err != nil {
        return err
    }

    if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram returned status %d", resp.StatusCode)
	}

    return nil
}
