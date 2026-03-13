package yookassa

import (
	"net/http"
	"testing"
)

func TestBasicAuthAuthorize(t *testing.T) {
	auth := BasicAuth{AccountID: "shop123", SecretKey: "secret456"}
	req, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)

	auth.Authorize(req)

	user, pass, ok := req.BasicAuth()
	if !ok {
		t.Fatal("expected basic auth to be set")
	}
	if user != "shop123" {
		t.Errorf("expected AccountID 'shop123', got '%s'", user)
	}
	if pass != "secret456" {
		t.Errorf("expected SecretKey 'secret456', got '%s'", pass)
	}
}

func TestBearerTokenAuthAuthorize(t *testing.T) {
	auth := BearerTokenAuth{Token: "my-oauth-token"}
	req, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)

	auth.Authorize(req)

	header := req.Header.Get("Authorization")
	expected := "Bearer my-oauth-token"
	if header != expected {
		t.Errorf("expected '%s', got '%s'", expected, header)
	}
}

func TestNewClientUsesBasicAuth(t *testing.T) {
	client := NewClient("shop123", "secret456")

	auth, ok := client.auth.(BasicAuth)
	if !ok {
		t.Fatal("expected BasicAuth authorizer")
	}
	if auth.AccountID != "shop123" || auth.SecretKey != "secret456" {
		t.Errorf("unexpected BasicAuth values: %+v", auth)
	}
}

func TestNewClientWithTokenUsesBearerAuth(t *testing.T) {
	client := NewClientWithToken("my-oauth-token")

	auth, ok := client.auth.(BearerTokenAuth)
	if !ok {
		t.Fatal("expected BearerTokenAuth authorizer")
	}
	if auth.Token != "my-oauth-token" {
		t.Errorf("expected token 'my-oauth-token', got '%s'", auth.Token)
	}
}
