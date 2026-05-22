package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestBearerAuth_ValidToken(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer valid_token_12345")

	middleware := BearerAuth()
	middleware(c)

	if c.IsAborted() {
		t.Error("middleware aborted unexpectedly")
	}

	token, exists := c.Get("bearer_token")
	if !exists {
		t.Fatal("bearer_token not set in context")
	}

	if token != "valid_token_12345" {
		t.Errorf("bearer_token = %v, want %q", token, "valid_token_12345")
	}
}

func TestBearerAuth_MissingToken(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	middleware := BearerAuth()
	middleware(c)

	if !c.IsAborted() {
		t.Error("middleware should have aborted")
	}

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusUnauthorized)
	}

	_, exists := c.Get("bearer_token")
	if exists {
		t.Error("bearer_token should not be set")
	}
}

func TestBearerAuth_InvalidFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		header string
	}{
		{"NoBearer", "token12345"},
		{"WrongPrefix", "Token abc123"},
		{"Basic", "Basic dXNlcjpwYXNz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
			c.Request.Header.Set("Authorization", tt.header)

			middleware := BearerAuth()
			middleware(c)

			if !c.IsAborted() {
				t.Error("middleware should have aborted for invalid format")
			}

			if w.Code != http.StatusUnauthorized {
				t.Errorf("status code = %d, want %d", w.Code, http.StatusUnauthorized)
			}
		})
	}
}

func TestBearerAuth_EmptyToken(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer ")

	middleware := BearerAuth()
	middleware(c)

	if !c.IsAborted() {
		t.Error("middleware should have aborted for empty token")
	}

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAPIKeyAuth_ValidKey(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Request.Header.Set("X-API-Key", "my_secret_key")

	middleware := APIKeyAuth()
	middleware(c)

	if c.IsAborted() {
		t.Error("middleware aborted unexpectedly")
	}

	apiKey, exists := c.Get("api_key")
	if !exists {
		t.Fatal("api_key not set in context")
	}

	if apiKey != "my_secret_key" {
		t.Errorf("api_key = %v, want %q", apiKey, "my_secret_key")
	}
}

func TestAPIKeyAuth_MissingKey(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	middleware := APIKeyAuth()
	middleware(c)

	if !c.IsAborted() {
		t.Error("middleware should have aborted")
	}

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusUnauthorized)
	}

	_, exists := c.Get("api_key")
	if exists {
		t.Error("api_key should not be set")
	}
}

func TestBasicAuth_ValidCredentials(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Request.Header.Set("Authorization", "Basic dXNlcjpwYXNzd29yZA==")

	middleware := BasicAuth()
	middleware(c)

	if c.IsAborted() {
		t.Error("middleware aborted unexpectedly")
	}
}

func TestBasicAuth_MissingAuth(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	middleware := BasicAuth()
	middleware(c)

	if !c.IsAborted() {
		t.Error("middleware should have aborted")
	}

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestBasicAuth_InvalidFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		header string
	}{
		{"NoBasic", "dXNlcjpwYXNz"},
		{"Bearer", "Bearer token123"},
		{"WrongPrefix", "Digest abc123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
			c.Request.Header.Set("Authorization", tt.header)

			middleware := BasicAuth()
			middleware(c)

			if !c.IsAborted() {
				t.Error("middleware should have aborted for invalid format")
			}

			if w.Code != http.StatusUnauthorized {
				t.Errorf("status code = %d, want %d", w.Code, http.StatusUnauthorized)
			}
		})
	}
}

func TestCORS_HeadersSet(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	middleware := CORS()
	middleware(c)

	if c.IsAborted() {
		t.Error("middleware aborted unexpectedly for GET request")
	}

	headers := w.Header()

	expectedHeaders := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD",
		"Access-Control-Allow-Headers": "Origin, Content-Type, Authorization, X-API-Key, X-Custom-Header",
		"Access-Control-Max-Age":       "86400",
	}

	for key, want := range expectedHeaders {
		got := headers.Get(key)
		if got != want {
			t.Errorf("header %q = %q, want %q", key, got, want)
		}
	}
}

func TestCORS_OptionsRequest(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodOptions, "/test", nil)

	middleware := CORS()
	middleware(c)

	if !c.IsAborted() {
		t.Error("middleware should abort for OPTIONS request")
	}

	if w.Code != http.StatusNoContent {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusNoContent)
	}

	headers := w.Header()
	if headers.Get("Access-Control-Allow-Origin") != "*" {
		t.Error("CORS headers not set for OPTIONS request")
	}
}

func TestCORS_NonOptionsRequest(t *testing.T) {
	t.Parallel()

	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodHead,
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(method, "/test", nil)

			middleware := CORS()
			middleware(c)

			if c.IsAborted() {
				t.Errorf("middleware should not abort for %s request", method)
			}

			if w.Header().Get("Access-Control-Allow-Origin") != "*" {
				t.Error("CORS headers not set")
			}
		})
	}
}

func TestRequestLogger_LogsRequest(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/test", nil)
	c.Request.Header.Set("Authorization", "Bearer test_token")
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("X-API-Key", "api_key_123")

	middleware := RequestLogger()

	c.Next()

	middleware(c)

	if c.IsAborted() {
		t.Error("middleware should not abort")
	}
}

func TestRequestLogger_HandlesLongHeaders(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/test", nil)

	longToken := "Bearer " + string(make([]byte, 100))
	c.Request.Header.Set("Authorization", longToken)

	middleware := RequestLogger()
	middleware(c)

	if c.IsAborted() {
		t.Error("middleware should not abort for long headers")
	}
}

func TestTruncate_ShortString(t *testing.T) {
	t.Parallel()

	input := "short"
	result := truncate(input, 10)

	if result != input {
		t.Errorf("truncate(%q, 10) = %q, want %q", input, result, input)
	}
}

func TestTruncate_ExactLength(t *testing.T) {
	t.Parallel()

	input := "exact"
	result := truncate(input, 5)

	if result != input {
		t.Errorf("truncate(%q, 5) = %q, want %q", input, result, input)
	}
}

func TestTruncate_LongString(t *testing.T) {
	t.Parallel()

	input := "this is a very long string"
	result := truncate(input, 10)

	expected := "this is a ..."
	if result != expected {
		t.Errorf("truncate(%q, 10) = %q, want %q", input, result, expected)
	}

	if len(result) != 13 {
		t.Errorf("truncated length = %d, want 13", len(result))
	}
}

func TestTruncate_EmptyString(t *testing.T) {
	t.Parallel()

	input := ""
	result := truncate(input, 10)

	if result != "" {
		t.Errorf("truncate(%q, 10) = %q, want empty string", input, result)
	}
}

func TestTruncate_ZeroMaxLen(t *testing.T) {
	t.Parallel()

	input := "test"
	result := truncate(input, 0)

	expected := "..."
	if result != expected {
		t.Errorf("truncate(%q, 0) = %q, want %q", input, result, expected)
	}
}

func TestMiddleware_ChainExecution(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)

	executed := false
	handler := func(c *gin.Context) {
		executed = true
		c.String(http.StatusOK, "ok")
	}

	router.Use(CORS())
	router.Use(BearerAuth())
	router.GET("/test", handler)

	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer valid_token")

	router.ServeHTTP(w, c.Request)

	if !executed {
		t.Error("handler not executed")
	}

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
	}

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("CORS header not set in chain")
	}
}

func TestMiddleware_ChainAbortOnAuth(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)

	executed := false
	handler := func(c *gin.Context) {
		executed = true
		c.String(http.StatusOK, "ok")
	}

	router.Use(CORS())
	router.Use(BearerAuth())
	router.GET("/test", handler)

	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	router.ServeHTTP(w, c.Request)

	if executed {
		t.Error("handler should not execute when auth fails")
	}

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusUnauthorized)
	}

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("CORS should still be set even when auth fails")
	}
}
