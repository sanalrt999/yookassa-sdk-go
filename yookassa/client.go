// Package yookassa implements all the necessary methods for working with YooMoney.
package yookassa

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

const (
	BaseURL = "https://api.yookassa.ru/v3/"
)

// Authorizer sets authorization headers on HTTP requests.
type Authorizer interface {
	Authorize(req *http.Request)
}

// BasicAuth implements Authorizer using HTTP Basic Authentication.
type BasicAuth struct {
	AccountID string
	SecretKey string
}

func (a BasicAuth) Authorize(req *http.Request) {
	req.SetBasicAuth(a.AccountID, a.SecretKey)
}

// BearerTokenAuth implements Authorizer using Bearer Token authentication.
type BearerTokenAuth struct {
	Token string
}

func (a BearerTokenAuth) Authorize(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+a.Token)
}

// Client works with YooMoney API.
type Client struct {
	client http.Client
	auth   Authorizer
}

func NewClient(accountId string, secretKey string) *Client {
	return &Client{
		client: http.Client{},
		auth:   BasicAuth{AccountID: accountId, SecretKey: secretKey},
	}
}

func NewClientWithToken(token string) *Client {
	return &Client{
		client: http.Client{},
		auth:   BearerTokenAuth{Token: token},
	}
}

func (c *Client) makeRequest(
	ctx context.Context,
	method string,
	endpoint string,
	body []byte,
	params map[string]interface{},
	idempotencyKey string,
) (*http.Response, error) {
	uri := fmt.Sprintf("%s%s", BaseURL, endpoint)

	req, err := http.NewRequestWithContext(ctx, method, uri, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	if idempotencyKey == "" {
		idempotencyKey = uuid.NewString()
	}

	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotence-Key", idempotencyKey)
	}

	c.auth.Authorize(req)

	if params != nil {
		q := req.URL.Query()
		for paramName, paramVal := range params {
			q.Add(paramName, fmt.Sprintf("%v", paramVal))
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
