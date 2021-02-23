package jwtheaders_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	jwtheaders "github.com/lion7/traefik-jwt-headers-plugin"
)

func TestJwt(t *testing.T) {
	/* #nosec */
	const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0Ijox" +
		"NTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	cfg := jwtheaders.CreateConfig()
	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := jwtheaders.New(ctx, next, cfg, "jwt-headers-plugin")
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

	// assertHeader(t, req, "X-JWT", "{\"sub\":\"1234567890\",\"name\":\"John Doe\",\"iat\":1516239022}")
	assertHeader(t, req, "X-JWT-sub", "1234567890")
	assertHeader(t, req, "X-JWT-name", "John Doe")
	assertHeader(t, req, "X-JWT-iat", "1516239022")
}

func TestNestedJwt(t *testing.T) {
	/* #nosec */
	const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0Ijox" +
		"NTE2MjM5MDIyLCJvcmdhbml6YXRpb24iOnsiaWQiOiIwOTg3NjU0MzIxIiwibmFtZSI6IkRvZSBjb21wYW55In19.uhtuQtJgnt_V9vsTsr" +
		"L9xoyYH8yOQYYG9KEGYjQT_zc"

	cfg := jwtheaders.CreateConfig()
	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := jwtheaders.New(ctx, next, cfg, "jwt-headers-plugin")
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

	// assertHeader(t, req, "X-JWT", "{\"sub\":\"1234567890\",\"name\":\"John Doe\",\"iat\":1516239022,"+
	//	"\"organization\":{\"id\":\"0987654321\",\"name\":\"Doe company\"}}")
	assertHeader(t, req, "X-JWT-sub", "1234567890")
	assertHeader(t, req, "X-JWT-name", "John Doe")
	assertHeader(t, req, "X-JWT-iat", "1516239022")
	assertHeader(t, req, "X-JWT-organization-id", "0987654321")
	assertHeader(t, req, "X-JWT-organization-name", "Doe company")
}

func TestNestedJwtWithConfig(t *testing.T) {
	/* #nosec */
	const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0Ijox" +
		"NTE2MjM5MDIyLCJvcmdhbml6YXRpb24iOnsiaWQiOiIwOTg3NjU0MzIxIiwibmFtZSI6IkRvZSBjb21wYW55In19.uhtuQtJgnt_V9vsTsr" +
		"L9xoyYH8yOQYYG9KEGYjQT_zc"

	cfg := jwtheaders.CreateConfig()
	cfg.DefaultMode = "drop"
	cfg.Claims["sub"] = "keep"
	cfg.Claims["organization.name"] = "keep"

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := jwtheaders.New(ctx, next, cfg, "jwt-headers-plugin")
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

	// assertHeader(t, req, "X-JWT", "{\"sub\":\"1234567890\",\"name\":\"John Doe\",\"iat\":1516239022,"+
	//	"\"organization\":{\"id\":\"0987654321\",\"name\":\"Doe company\"}}")
	assertHeader(t, req, "X-JWT-sub", "1234567890")
	assertHeader(t, req, "X-JWT-name", "")
	assertHeader(t, req, "X-JWT-iat", "")
	assertHeader(t, req, "X-JWT-organization-id", "")
	assertHeader(t, req, "X-JWT-organization-name", "Doe company")
}

func assertHeader(t *testing.T, req *http.Request, key, expected string) {
	t.Helper()

	if req.Header.Get(key) != expected {
		t.Errorf("invalid header value: %s", req.Header.Get(key))
	}
}
