package jwtlogging

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
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// JwtLogging a JwtLogging plugin.
type JwtLogging struct {
	next     http.Handler
	name     string
	template *template.Template
}

// New created a new JwtLogging plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &JwtLogging{
		next:     next,
		name:     name,
		template: template.New("jwt-logging").Delims("[[", "]]"),
	}, nil
}

func (a *JwtLogging) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, value := range req.Header.Values("Authorization") {
		if strings.HasPrefix(value, "Bearer ") {
			token := strings.TrimPrefix(value, "Bearer ")
			body := strings.Split(token, ".")[1]
			decodedBody, err := base64.RawStdEncoding.DecodeString(body)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			jsonBody := string(decodedBody)
			req.Header.Set("X-JWT", jsonBody)

			jsonMap := make(map[string]interface{})
			dec := json.NewDecoder(bytes.NewReader(decodedBody))
			dec.UseNumber()
			err = dec.Decode(&jsonMap)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			for key, value := range jsonMap {
				s := fmt.Sprintf("%v", value)
				req.Header.Set("X-JWT-"+key, s)
			}
		}
	}

	a.next.ServeHTTP(rw, req)
}
