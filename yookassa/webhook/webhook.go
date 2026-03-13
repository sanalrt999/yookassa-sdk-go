package yoowebhook

// Webhook represents a webhook subscription.
type Webhook struct {
	ID    string           `json:"id"`
	Event WebhookEventType `json:"event"`
	URL   string           `json:"url"`
}

// WebhookList represents a list of webhook subscriptions.
type WebhookList struct {
	Type  string    `json:"type"`
	Items []Webhook `json:"items"`
}

// CreateWebhookRequest is the request to create a webhook subscription.
type CreateWebhookRequest struct {
	Event WebhookEventType `json:"event"`
	URL   string           `json:"url"`
}
