package traefik_jwt_logging_plugin_test

import (
	"context"
	traefik_jwt_logging_plugin "github.com/lion7/traefik-jwt-logging-plugin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJwtLogging(t *testing.T) {
	cfg := traefik_jwt_logging_plugin.CreateConfig()
	cfg.Headers["X-Host"] = "[[.Host]]"
	cfg.Headers["X-Method"] = "[[.Method]]"
	cfg.Headers["X-URL"] = "[[.URL]]"
	cfg.Headers["X-URL"] = "[[.URL]]"
	cfg.Headers["X-JwtLogging"] = "test"

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := traefik_jwt_logging_plugin.New(ctx, next, cfg, "jwt-logging-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	assertHeader(t, req, "X-Host", "localhost")
	assertHeader(t, req, "X-URL", "http://localhost")
	assertHeader(t, req, "X-Method", "GET")
	assertHeader(t, req, "X-JwtLogging", "test")
}

func assertHeader(t *testing.T, req *http.Request, key, expected string) {
	t.Helper()

	if req.Header.Get(key) != expected {
		t.Errorf("invalid header value: %s", req.Header.Get(key))
	}
}
