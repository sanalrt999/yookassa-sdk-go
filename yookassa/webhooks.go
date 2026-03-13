package yookassa

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	yooerror "github.com/sanalrt999/yookassa-sdk-go/yookassa/errors"
	yoowebhook "github.com/sanalrt999/yookassa-sdk-go/yookassa/webhook"
)

const (
	WebhooksEndpoint = "webhooks"
)

// WebhookHandler works with webhook subscriptions.
type WebhookHandler struct {
	client         *Client
	idempotencyKey string
}

func NewWebhookHandler(client *Client) *WebhookHandler {
	return &WebhookHandler{client: client}
}

func (w WebhookHandler) WithIdempotencyKey(idempotencyKey string) WebhookHandler {
	w.idempotencyKey = idempotencyKey

	return w
}

// CreateWebhookCtx creates a new webhook subscription.
func (w *WebhookHandler) CreateWebhookCtx(ctx context.Context, req *yoowebhook.CreateWebhookRequest) (*yoowebhook.Webhook, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := w.client.makeRequest(ctx, http.MethodPost, WebhooksEndpoint, body, nil, w.idempotencyKey)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var respError error
		respError, err = yooerror.GetError(resp.Body)
		if err != nil {
			return nil, err
		}

		return nil, respError
	}

	var responseBytes []byte
	responseBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	webhook := yoowebhook.Webhook{}
	err = json.Unmarshal(responseBytes, &webhook)
	if err != nil {
		return nil, err
	}
	return &webhook, nil
}

// CreateWebhook creates a new webhook subscription.
// Deprecated: Use CreateWebhookCtx instead.
func (w *WebhookHandler) CreateWebhook(req *yoowebhook.CreateWebhookRequest) (*yoowebhook.Webhook, error) {
	return w.CreateWebhookCtx(context.Background(), req)
}

// ListWebhooksCtx returns the list of webhook subscriptions.
func (w *WebhookHandler) ListWebhooksCtx(ctx context.Context) (*yoowebhook.WebhookList, error) {
	resp, err := w.client.makeRequest(ctx, http.MethodGet, WebhooksEndpoint, nil, nil, w.idempotencyKey)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var respError error
		respError, err = yooerror.GetError(resp.Body)
		if err != nil {
			return nil, err
		}

		return nil, respError
	}

	var responseBytes []byte
	responseBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	webhookList := yoowebhook.WebhookList{}
	err = json.Unmarshal(responseBytes, &webhookList)
	if err != nil {
		return nil, err
	}
	return &webhookList, nil
}

// ListWebhooks returns the list of webhook subscriptions.
// Deprecated: Use ListWebhooksCtx instead.
func (w *WebhookHandler) ListWebhooks() (*yoowebhook.WebhookList, error) {
	return w.ListWebhooksCtx(context.Background())
}

// DeleteWebhookCtx deletes a webhook subscription by ID.
func (w *WebhookHandler) DeleteWebhookCtx(ctx context.Context, webhookID string) error {
	endpoint := fmt.Sprintf("%s/%s", WebhooksEndpoint, webhookID)

	resp, err := w.client.makeRequest(ctx, http.MethodDelete, endpoint, nil, nil, w.idempotencyKey)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var respError error
		respError, err = yooerror.GetError(resp.Body)
		if err != nil {
			return err
		}

		return respError
	}

	return nil
}

// DeleteWebhook deletes a webhook subscription by ID.
// Deprecated: Use DeleteWebhookCtx instead.
func (w *WebhookHandler) DeleteWebhook(webhookID string) error {
	return w.DeleteWebhookCtx(context.Background(), webhookID)
}
