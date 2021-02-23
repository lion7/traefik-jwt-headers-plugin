package jwtlogging_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	jwtlogging "github.com/lion7/traefik-jwt-logging-plugin"
)

func TestJwt(t *testing.T) {
	/* #nosec */
	const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0Ijox" +
		"NTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	cfg := jwtlogging.CreateConfig()
	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := jwtlogging.New(ctx, next, cfg, "jwt-logging-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	handler.ServeHTTP(recorder, req)

	assertHeader(t, req, "X-JWT", "{\"sub\":\"1234567890\",\"name\":\"John Doe\",\"iat\":1516239022}")
	assertHeader(t, req, "X-JWT-sub", "1234567890")
	assertHeader(t, req, "X-JWT-name", "John Doe")
	assertHeader(t, req, "X-JWT-iat", "1516239022")
}

func TestNestedJwt(t *testing.T) {
	/* #nosec */
	const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0Ijox" +
		"NTE2MjM5MDIyLCJvcmdhbml6YXRpb24iOnsiaWQiOiIwOTg3NjU0MzIxIiwibmFtZSI6IkRvZSBjb21wYW55In19.uhtuQtJgnt_V9vsTsr" +
		"L9xoyYH8yOQYYG9KEGYjQT_zc"

	cfg := jwtlogging.CreateConfig()
	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := jwtlogging.New(ctx, next, cfg, "jwt-logging-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	handler.ServeHTTP(recorder, req)

	assertHeader(t, req, "X-JWT", "{\"sub\":\"1234567890\",\"name\":\"John Doe\",\"iat\":1516239022}")
	assertHeader(t, req, "X-JWT-sub", "1234567890")
	assertHeader(t, req, "X-JWT-name", "John Doe")
	assertHeader(t, req, "X-JWT-iat", "1516239022")
}

func assertHeader(t *testing.T, req *http.Request, key, expected string) {
	t.Helper()

	if req.Header.Get(key) != expected {
		t.Errorf("invalid header value: %s", req.Header.Get(key))
	}
}
