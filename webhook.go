package gchat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type (
	PostWebhookOptions struct {
		Endpoint string
		Payload  WebhookPayload
	}
	WebhookPayload struct {
		Text string `json:"text"`
	}
)

type WebhookClient struct {
	Endpoint string
	Client   *http.Client
}

func (c *WebhookClient) Post(ctx context.Context, p WebhookPayload) error {
	b, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("gchat: marshal webhook payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.Endpoint, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("gchat: create http request: %w", err)
	}
	res, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("gchat: post message: %w", err)
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("gchat: post message: %v", res.Status)
	}
	defer res.Body.Close()
	_, err = io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("gchat: read response: %w", err)
	}
	return nil
}
