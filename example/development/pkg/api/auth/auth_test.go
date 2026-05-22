package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"example/development/pkg/middleware"
	"example/development/pkg/model"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestServer(t *testing.T) (*gin.Engine, huma.API) {
	t.Helper()

	r := gin.New()
	r.Use(middleware.CORS())

	config := huma.DefaultConfig("Test API", "1.0.0")
	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearerAuth": {
			Type:   "http",
			Scheme: "bearer",
		},
		"apiKeyAuth": {
			Type: "apiKey",
			Name: "X-API-Key",
			In:   "header",
		},
		"basicAuth": {
			Type:   "http",
			Scheme: "basic",
		},
	}

	api := humagin.New(r, config)
	return r, api
}

func TestMessages_Structure(t *testing.T) {
	t.Parallel()

	msg := Messages{
		PublicMessage: "Public access",
		BearerSuccess: "Bearer success",
		APIKeySuccess: "API Key success",
		BasicSuccess:  "Basic success",
		PublicSummary: "Public summary",
		BearerSummary: "Bearer summary",
		APIKeySummary: "API Key summary",
		BasicSummary:  "Basic summary",
		Tag:           "Authentication",
	}

	if msg.PublicMessage != "Public access" {
		t.Errorf("PublicMessage = %q, want %q", msg.PublicMessage, "Public access")
	}

	if msg.BearerSuccess != "Bearer success" {
		t.Errorf("BearerSuccess = %q, want %q", msg.BearerSuccess, "Bearer success")
	}

	if msg.APIKeySuccess != "API Key success" {
		t.Errorf("APIKeySuccess = %q, want %q", msg.APIKeySuccess, "API Key success")
	}

	if msg.BasicSuccess != "Basic success" {
		t.Errorf("BasicSuccess = %q, want %q", msg.BasicSuccess, "Basic success")
	}

	if msg.Tag != "Authentication" {
		t.Errorf("Tag = %q, want %q", msg.Tag, "Authentication")
	}
}

func TestOutput_Structure(t *testing.T) {
	t.Parallel()

	output := Output{
		Body: model.SuccessResponse(map[string]string{
			"message": "test",
		}),
	}

	if output.Body.Code != 200 {
		t.Errorf("Code = %d, want 200", output.Body.Code)
	}

	if output.Body.Data == nil {
		t.Error("Data is nil")
	}
}

func TestRegister_DoesNotPanic(t *testing.T) {
	t.Parallel()

	_, api := humatest.New(t)

	msg := Messages{
		PublicMessage: "Public access",
		BearerSuccess: "Bearer OK",
		APIKeySuccess: "API Key OK",
		BasicSuccess:  "Basic OK",
		PublicSummary: "Public endpoint",
		BearerSummary: "Bearer endpoint",
		APIKeySummary: "API Key endpoint",
		BasicSummary:  "Basic endpoint",
		Tag:           "Auth",
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Register panicked: %v", r)
		}
	}()

	Register(api, msg)
}

func TestMessages_AllFieldsSet(t *testing.T) {
	t.Parallel()

	msg := Messages{
		PublicMessage: "msg1",
		BearerSuccess: "msg2",
		APIKeySuccess: "msg3",
		BasicSuccess:  "msg4",
		PublicSummary: "sum1",
		BearerSummary: "sum2",
		APIKeySummary: "sum3",
		BasicSummary:  "sum4",
		Tag:           "tag",
	}

	fields := []string{
		msg.PublicMessage,
		msg.BearerSuccess,
		msg.APIKeySuccess,
		msg.BasicSuccess,
		msg.PublicSummary,
		msg.BearerSummary,
		msg.APIKeySummary,
		msg.BasicSummary,
		msg.Tag,
	}

	for i, field := range fields {
		if field == "" {
			t.Errorf("field %d is empty", i)
		}
	}
}

// HTTP Integration Tests

func TestAuthPublic_HTTP(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		PublicMessage: "Public access granted",
		PublicSummary: "Public endpoint",
		Tag:           "Auth",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodGet, "/auth/public", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
	}

	var resp model.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Code != 200 {
		t.Errorf("response code = %d, want 200", resp.Code)
	}
}

func TestAuthBearer_HTTP_Success(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		BearerSuccess: "Bearer token accepted",
		BearerSummary: "Bearer endpoint",
		Tag:           "Auth",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodGet, "/auth/bearer", nil)
	req.Header.Set("Authorization", "Bearer test_token_12345")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
		t.Logf("Response body: %s", w.Body.String())
	}

	var resp model.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Code != 200 {
		t.Errorf("response code = %d, want 200", resp.Code)
	}
}

func TestAuthBearer_HTTP_WithDifferentTokens(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		BearerSuccess: "Bearer token accepted",
		BearerSummary: "Bearer endpoint",
		Tag:           "Auth",
	}
	Register(api, msg)

	tests := []struct {
		name  string
		token string
	}{
		{"Valid Bearer token", "Bearer valid_token_123"},
		{"Another valid token", "Bearer another_token"},
		{"Token with special chars", "Bearer token-with_special.chars"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/auth/bearer", nil)
			req.Header.Set("Authorization", tt.token)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("status = %d, want 200", w.Code)
			}
		})
	}
}

func TestAuthAPIKey_HTTP_Success(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		APIKeySuccess: "API Key accepted",
		APIKeySummary: "API Key endpoint",
		Tag:           "Auth",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodGet, "/auth/api-key", nil)
	req.Header.Set("X-API-Key", "my_secret_key_123")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
		t.Logf("Response body: %s", w.Body.String())
	}

	var resp model.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Code != 200 {
		t.Errorf("response code = %d, want 200", resp.Code)
	}
}

func TestAuthAPIKey_HTTP_WithDifferentKeys(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		APIKeySuccess: "API Key accepted",
		APIKeySummary: "API Key endpoint",
		Tag:           "Auth",
	}
	Register(api, msg)

	tests := []struct {
		name   string
		apiKey string
	}{
		{"Standard key", "my_api_key_123"},
		{"UUID key", "550e8400-e29b-41d4-a716-446655440000"},
		{"Complex key", "sk_live_51H...abc123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/auth/api-key", nil)
			req.Header.Set("X-API-Key", tt.apiKey)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("status = %d, want 200", w.Code)
			}
		})
	}
}

func TestAuthBasic_HTTP_Success(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		BasicSuccess: "Basic auth accepted",
		BasicSummary: "Basic endpoint",
		Tag:          "Auth",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodGet, "/auth/basic", nil)
	req.Header.Set("Authorization", "Basic dXNlcjpwYXNzd29yZA==")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
		t.Logf("Response body: %s", w.Body.String())
	}

	var resp model.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Code != 200 {
		t.Errorf("response code = %d, want 200", resp.Code)
	}
}

func TestAuthBasic_HTTP_WithDifferentCredentials(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		BasicSuccess: "Basic auth accepted",
		BasicSummary: "Basic endpoint",
		Tag:          "Auth",
	}
	Register(api, msg)

	tests := []struct {
		name   string
		header string
	}{
		{"Standard credentials", "Basic dXNlcjpwYXNzd29yZA=="},
		{"Admin credentials", "Basic YWRtaW46YWRtaW4xMjM="},
		{"Email as username", "Basic dGVzdEB0ZXN0LmNvbTpwYXNzd29yZA=="},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/auth/basic", nil)
			req.Header.Set("Authorization", tt.header)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("status = %d, want 200", w.Code)
			}
		})
	}
}

func TestAuth_ResponseFormat(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		PublicMessage: "Test message",
		PublicSummary: "Public endpoint",
		Tag:           "Auth",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodGet, "/auth/public", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", w.Code, http.StatusOK)
	}

	var resp model.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Code != 200 {
		t.Errorf("code = %d, want 200", resp.Code)
	}

	if resp.Data == nil {
		t.Error("data is nil")
	}
}

func TestAuth_AllEndpoints(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		PublicMessage: "Public",
		BearerSuccess: "Bearer",
		APIKeySuccess: "API Key",
		BasicSuccess:  "Basic",
		PublicSummary: "Public",
		BearerSummary: "Bearer",
		APIKeySummary: "API Key",
		BasicSummary:  "Basic",
		Tag:           "Auth",
	}
	Register(api, msg)

	tests := []struct {
		name       string
		path       string
		headers    map[string]string
		wantStatus int
	}{
		{
			name:       "Public endpoint",
			path:       "/auth/public",
			wantStatus: http.StatusOK,
		},
		{
			name: "Bearer with valid token",
			path: "/auth/bearer",
			headers: map[string]string{
				"Authorization": "Bearer token123",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "API Key with valid key",
			path: "/auth/api-key",
			headers: map[string]string{
				"X-API-Key": "key123",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Basic with credentials",
			path: "/auth/basic",
			headers: map[string]string{
				"Authorization": "Basic dXNlcjpwYXNz",
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}
