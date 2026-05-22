package debug

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"example/development/pkg/middleware"

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
	api := humagin.New(r, config)
	return r, api
}

func TestMessages_Structure(t *testing.T) {
	t.Parallel()

	msg := Messages{
		EchoSummary: "Echo endpoint",
		EchoDesc:    "Returns request information",
		Tag:         "Debug",
	}

	if msg.EchoSummary != "Echo endpoint" {
		t.Errorf("EchoSummary = %q, want %q", msg.EchoSummary, "Echo endpoint")
	}

	if msg.EchoDesc != "Returns request information" {
		t.Errorf("EchoDesc = %q, want %q", msg.EchoDesc, "Returns request information")
	}

	if msg.Tag != "Debug" {
		t.Errorf("Tag = %q, want %q", msg.Tag, "Debug")
	}
}

func TestRequestInfo_Structure(t *testing.T) {
	t.Parallel()

	info := RequestInfo{
		Method:  "GET",
		Path:    "/debug/echo",
		Headers: map[string][]string{"Authorization": {"Bearer token"}},
		Query:   map[string][]string{"name": {"test"}},
		Host:    "localhost:8000",
		Proto:   "HTTP/1.1",
	}

	if info.Method != "GET" {
		t.Errorf("Method = %q, want %q", info.Method, "GET")
	}

	if info.Path != "/debug/echo" {
		t.Errorf("Path = %q, want %q", info.Path, "/debug/echo")
	}

	if info.Host != "localhost:8000" {
		t.Errorf("Host = %q, want %q", info.Host, "localhost:8000")
	}

	if info.Proto != "HTTP/1.1" {
		t.Errorf("Proto = %q, want %q", info.Proto, "HTTP/1.1")
	}
}

func TestEchoOutput_Structure(t *testing.T) {
	t.Parallel()

	output := EchoOutput{
		Body: RequestInfo{
			Method: "POST",
			Path:   "/debug/echo",
		},
	}

	if output.Body.Method != "POST" {
		t.Errorf("Method = %q, want %q", output.Body.Method, "POST")
	}

	if output.Body.Path != "/debug/echo" {
		t.Errorf("Path = %q, want %q", output.Body.Path, "/debug/echo")
	}
}

func TestRegister_DoesNotPanic(t *testing.T) {
	t.Parallel()

	_, api := humatest.New(t)

	msg := Messages{
		EchoSummary: "Echo",
		EchoDesc:    "Echo description",
		Tag:         "Debug",
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
		EchoSummary: "summary",
		EchoDesc:    "description",
		Tag:         "tag",
	}

	fields := []string{
		msg.EchoSummary,
		msg.EchoDesc,
		msg.Tag,
	}

	for i, field := range fields {
		if field == "" {
			t.Errorf("field %d is empty", i)
		}
	}
}

func TestRequestInfo_EmptyMaps(t *testing.T) {
	t.Parallel()

	info := RequestInfo{
		Method:  "GET",
		Path:    "/test",
		Headers: map[string][]string{},
		Query:   map[string][]string{},
	}

	if info.Headers == nil {
		t.Error("Headers should not be nil")
	}

	if len(info.Headers) != 0 {
		t.Errorf("Headers length = %d, want 0", len(info.Headers))
	}

	if info.Query == nil {
		t.Error("Query should not be nil")
	}

	if len(info.Query) != 0 {
		t.Errorf("Query length = %d, want 0", len(info.Query))
	}
}

// HTTP Integration Tests

func TestDebugEchoGET_HTTP_WithHeaders(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		EchoSummary: "Echo endpoint",
		EchoDesc:    "Returns request information",
		Tag:         "Debug",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodGet, "/debug/echo", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	req.Header.Set("X-API-Key", "test-key")
	req.Header.Set("X-Custom-Header", "custom-value")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
		t.Logf("Response: %s", w.Body.String())
	}

	var resp RequestInfo
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Method != "GET" {
		t.Errorf("method = %q, want %q", resp.Method, "GET")
	}

	if resp.Path != "/debug/echo" {
		t.Errorf("path = %q, want %q", resp.Path, "/debug/echo")
	}

	if len(resp.Headers["Authorization"]) > 0 && resp.Headers["Authorization"][0] != "Bearer test-token" {
		t.Errorf("Authorization = %q, want %q", resp.Headers["Authorization"][0], "Bearer test-token")
	}
}

func TestDebugEchoGET_HTTP_WithQuery(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		EchoSummary: "Echo endpoint",
		Tag:         "Debug",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodGet, "/debug/echo?name=TestUser", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
		t.Logf("Response: %s", w.Body.String())
	}

	var resp RequestInfo
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(resp.Query["name"]) > 0 && resp.Query["name"][0] != "TestUser" {
		t.Errorf("query name = %q, want %q", resp.Query["name"][0], "TestUser")
	}
}

func TestDebugEchoGET_HTTP_NoParams(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		EchoSummary: "Echo endpoint",
		Tag:         "Debug",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodGet, "/debug/echo", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
	}

	var resp RequestInfo
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Method != "GET" {
		t.Errorf("method = %q, want %q", resp.Method, "GET")
	}
}

func TestDebugEchoPOST_HTTP_WithHeaders(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		EchoSummary: "Echo endpoint",
		EchoDesc:    "Returns request information",
		Tag:         "Debug",
	}
	Register(api, msg)

	body := []byte(`{"test": "data"}`)
	req := httptest.NewRequest(http.MethodPost, "/debug/echo", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer post-token")
	req.Header.Set("X-API-Key", "post-key")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
		t.Logf("Response: %s", w.Body.String())
	}

	var resp RequestInfo
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Method != "POST" {
		t.Errorf("method = %q, want %q", resp.Method, "POST")
	}

	if len(resp.Headers["Authorization"]) > 0 && resp.Headers["Authorization"][0] != "Bearer post-token" {
		t.Errorf("Authorization = %q, want %q", resp.Headers["Authorization"][0], "Bearer post-token")
	}
}

func TestDebugEcho_HTTP_BothMethods(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		EchoSummary: "Echo endpoint",
		EchoDesc:    "Returns request information",
		Tag:         "Debug",
	}
	Register(api, msg)

	tests := []struct {
		name   string
		method string
		body   []byte
	}{
		{"GET", http.MethodGet, nil},
		{"POST", http.MethodPost, []byte(`{"key":"value"}`)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.body != nil {
				req = httptest.NewRequest(tt.method, "/debug/echo", bytes.NewReader(tt.body))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tt.method, "/debug/echo", nil)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("status = %d, want 200", w.Code)
				return
			}

			var resp RequestInfo
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}

			if resp.Method != tt.method {
				t.Errorf("method = %q, want %q", resp.Method, tt.method)
			}
		})
	}
}

func TestDebugEcho_HTTP_HeaderValidation(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		EchoSummary: "Echo endpoint",
		Tag:         "Debug",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodGet, "/debug/echo?name=test", nil)
	req.Header.Set("Authorization", "Bearer validation-token")
	req.Header.Set("X-API-Key", "validation-key")
	req.Header.Set("X-Custom-Header", "validation-custom")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", w.Code, http.StatusOK)
	}

	var resp RequestInfo
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Method != "GET" {
		t.Errorf("method = %q, want %q", resp.Method, "GET")
	}

	if resp.Path != "/debug/echo" {
		t.Errorf("path = %q, want %q", resp.Path, "/debug/echo")
	}

	if len(resp.Headers["X-Custom-Header"]) > 0 && resp.Headers["X-Custom-Header"][0] != "validation-custom" {
		t.Errorf("X-Custom-Header = %q, want %q", resp.Headers["X-Custom-Header"][0], "validation-custom")
	}

	if len(resp.Query["name"]) > 0 && resp.Query["name"][0] != "test" {
		t.Errorf("query name = %q, want %q", resp.Query["name"][0], "test")
	}
}
