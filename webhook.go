// gchat is a simple client for posting text messages to Google Chat via webhook
package gchat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var ErrNoEndpoint = errors.New("gchat: no endpoint provided")

type WebhookPayload struct {
	Text string `json:"text"`
}

type WebhookClient struct {
	Endpoint string
	Client   *http.Client
}

func (c *WebhookClient) Post(ctx context.Context, p WebhookPayload) error {
	if c.Endpoint == "" {
		return ErrNoEndpoint
	}

	b, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("gchat: marshal webhook payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.Endpoint, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("gchat: create http request: %w", err)
	}

	client := c.Client
	if c.Client == nil {
		client = http.DefaultClient
	}

	res, err := client.Do(req)
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
