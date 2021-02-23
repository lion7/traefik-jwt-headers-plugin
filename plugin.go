// Package traefik_jwt_headers_plugin - Traefik middleware plugin which forwards JWT claims as request headers
package traefik_jwt_headers_plugin

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"text/template"
)

// Config the plugin configuration.
type Config struct {
	DefaultMode string            `json:"defaultMode"`
	Claims      map[string]string `json:"claims,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		DefaultMode: "keep",
		Claims:      make(map[string]string),
	}
}

// JwtHeaders a JwtHeaders plugin.
type JwtHeaders struct {
	next        http.Handler
	defaultMode string
	headers     map[string]string
	name        string
	template    *template.Template
}

// New created a new JwtHeaders plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &JwtHeaders{
		defaultMode: config.DefaultMode,
		headers:     config.Claims,
		next:        next,
		name:        name,
		template:    template.New("jwt-headers").Delims("[[", "]]"),
	}, nil
}

func (a *JwtHeaders) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, value := range req.Header.Values("Authorization") {
		if strings.HasPrefix(value, "Bearer ") {
			token := strings.TrimPrefix(value, "Bearer ")
			body := strings.Split(token, ".")[1]

			decodedBody, err := base64.RawStdEncoding.DecodeString(body)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			// jsonBody := string(decodedBody)
			// req.Header.Set("X-JWT", jsonBody)

			jsonMap := make(map[string]interface{})
			dec := json.NewDecoder(bytes.NewReader(decodedBody))
			dec.UseNumber()

			err = dec.Decode(&jsonMap)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			a.setHeaders(jsonMap, req, "")
		}
	}

	a.next.ServeHTTP(rw, req)
}

func (a *JwtHeaders) setHeaders(jsonMap map[string]interface{}, req *http.Request, path string) {
	if len(a.headers) == 0 && a.defaultMode != "keep" {
		return
	}

	for key, value := range jsonMap {
		newPath := ""
		if len(path) == 0 {
			newPath = key
		} else {
			newPath = path + "." + key
		}

		nestedMap, ok := value.(map[string]interface{})
		if ok {
			a.setHeaders(nestedMap, req, newPath)
		} else {
			mode, ok := a.headers[newPath]
			if !ok {
				mode = a.defaultMode
			}

			if mode == "keep" {
				headerKey := "X-JWT-" + strings.ReplaceAll(newPath, ".", "-")
				headerValue := fmt.Sprintf("%v", value)
				req.Header.Set(headerKey, headerValue)
			}
		}
	}
}
