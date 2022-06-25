package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/nna774/gyazo"
)

type APIConfig struct {
	SlackWebhook string
	GyazoToken   string
}
type Message struct {
	Text     string `json:"text"`
	UserName string `json:"username"`
	Emoji    string `json:"icon_emoji"`
}

func SendMessage(msg Message, config APIConfig) error {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal : %w", err)
	}
	resp, err := http.PostForm(config.SlackWebhook, url.Values{"payload": {string(bytes)}})
	if err != nil {
		return fmt.Errorf("post form : %w", err)
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return fmt.Errorf("faile slack post : %s", resp.Status)
	}
	return nil
}

func UploadImage(img io.Reader, config APIConfig) (string, error) {
	client := gyazo.NewOauth2Client(config.GyazoToken)
	result, err := client.Upload(img, nil)
	if err != nil {
		return "", fmt.Errorf("gyazo upload : %w", err)
	}
	return result.URL, nil
}
