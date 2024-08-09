package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
)

func calcHMACSHA1(message, secret string) string {
	mac := hmac.New(sha1.New, []byte(secret))
	_, _ = mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

func postTraqWebhook(webhookID string, secret string, text string) error {
	url := fmt.Sprintf("https://q.trap.jp/api/v3/webhooks/%s", webhookID)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(text))
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}

	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	req.Header.Set("X-Traq-Signature", calcHMACSHA1(text, secret))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to post data to traQ: %v", err)
	}

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to post data to traQ: %v", res.Status)
	}

	return nil
}
