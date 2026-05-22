package params

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
	api := humagin.New(r, config)
	return r, api
}

func TestMessages_Structure(t *testing.T) {
	t.Parallel()

	msg := Messages{
		PathSummary:   "Path parameters",
		QuerySummary:  "Query parameters",
		HeaderSummary: "Header parameters",
		CookieSummary: "Cookie parameters",
		MixedSummary:  "Mixed parameters",
		Tag:           "Parameters",
	}

	if msg.PathSummary != "Path parameters" {
		t.Errorf("PathSummary = %q, want %q", msg.PathSummary, "Path parameters")
	}

	if msg.QuerySummary != "Query parameters" {
		t.Errorf("QuerySummary = %q, want %q", msg.QuerySummary, "Query parameters")
	}

	if msg.HeaderSummary != "Header parameters" {
		t.Errorf("HeaderSummary = %q, want %q", msg.HeaderSummary, "Header parameters")
	}

	if msg.CookieSummary != "Cookie parameters" {
		t.Errorf("CookieSummary = %q, want %q", msg.CookieSummary, "Cookie parameters")
	}

	if msg.MixedSummary != "Mixed parameters" {
		t.Errorf("MixedSummary = %q, want %q", msg.MixedSummary, "Mixed parameters")
	}

	if msg.Tag != "Parameters" {
		t.Errorf("Tag = %q, want %q", msg.Tag, "Parameters")
	}
}

func TestOutput_Structure(t *testing.T) {
	t.Parallel()

	output := Output{
		Body: model.ParamTestResponse{
			PathParams:  map[string]string{"id": "123"},
			QueryParams: map[string][]string{"name": {"test"}},
			Headers:     map[string][]string{"Authorization": {"Bearer token"}},
			Cookies:     map[string]string{"session": "abc"},
		},
	}

	if output.Body.PathParams == nil {
		t.Error("PathParams is nil")
	}

	if output.Body.PathParams["id"] != "123" {
		t.Errorf("PathParams[id] = %q, want %q", output.Body.PathParams["id"], "123")
	}
}

func TestRegister_DoesNotPanic(t *testing.T) {
	t.Parallel()

	_, api := humatest.New(t)

	msg := Messages{
		PathSummary:   "Path",
		QuerySummary:  "Query",
		HeaderSummary: "Header",
		CookieSummary: "Cookie",
		MixedSummary:  "Mixed",
		Tag:           "Params",
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
		PathSummary:   "path",
		QuerySummary:  "query",
		HeaderSummary: "header",
		CookieSummary: "cookie",
		MixedSummary:  "mixed",
		Tag:           "tag",
	}

	fields := []string{
		msg.PathSummary,
		msg.QuerySummary,
		msg.HeaderSummary,
		msg.CookieSummary,
		msg.MixedSummary,
		msg.Tag,
	}

	for i, field := range fields {
		if field == "" {
			t.Errorf("field %d is empty", i)
		}
	}
}

func TestOutput_EmptyMaps(t *testing.T) {
	t.Parallel()

	output := Output{
		Body: model.ParamTestResponse{
			PathParams:  map[string]string{},
			QueryParams: map[string][]string{},
			Headers:     map[string][]string{},
			Cookies:     map[string]string{},
		},
	}

	if output.Body.PathParams == nil {
		t.Error("PathParams should not be nil")
	}

	if len(output.Body.PathParams) != 0 {
		t.Errorf("PathParams length = %d, want 0", len(output.Body.PathParams))
	}
}

// HTTP Integration Tests

func TestParamsPath_HTTP_ValidParams(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		PathSummary: "Path parameters",
		Tag:         "Params",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodGet, "/params/path/123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
		t.Logf("Response: %s", w.Body.String())
	}

	var resp model.ParamTestResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.PathParams["id"] != "123" {
		t.Errorf("path param id = %q, want %q", resp.PathParams["id"], "123")
	}
}

func TestParamsPath_HTTP_DifferentIDs(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		PathSummary: "Path parameters",
		Tag:         "Params",
	}
	Register(api, msg)

	tests := []struct {
		name string
		id   string
	}{
		{"Numeric ID", "12345"},
		{"UUID", "550e8400-e29b-41d4-a716-446655440000"},
		{"Alphanumeric", "abc123xyz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/params/path/"+tt.id, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("status = %d, want 200", w.Code)
			}

			var resp model.ParamTestResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}

			if resp.PathParams["id"] != tt.id {
				t.Errorf("id = %q, want %q", resp.PathParams["id"], tt.id)
			}
		})
	}
}

func TestParamsQuery_HTTP_SingleValue(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		QuerySummary: "Query parameters",
		Tag:          "Params",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodGet, "/params/query?name=John&age=30", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
		t.Logf("Response: %s", w.Body.String())
	}

	var resp model.ParamTestResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(resp.QueryParams["name"]) > 0 && resp.QueryParams["name"][0] != "John" {
		t.Errorf("query param name = %q, want %q", resp.QueryParams["name"][0], "John")
	}
}

func TestParamsQuery_HTTP_EmptyValue(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		QuerySummary: "Query parameters",
		Tag:          "Params",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodGet, "/params/query", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestParamsQuery_HTTP_SpecialChars(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		QuerySummary: "Query parameters",
		Tag:          "Params",
	}
	Register(api, msg)

	tests := []struct {
		name string
		url  string
	}{
		{"Spaces", "/params/query?name=John+Doe&age=25"},
		{"Special chars", "/params/query?name=user%2Btest%40example.com&age=25"},
		{"Unicode", "/params/query?name=%E7%94%A8%E6%88%B7%E5%90%8D&age=25"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("status = %d, want 200", w.Code)
			}
		})
	}
}

func TestParamsHeader_HTTP_ValidHeaders(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		HeaderSummary: "Header parameters",
		Tag:           "Params",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodGet, "/params/header", nil)
	req.Header.Set("User-Agent", "TestAgent/1.0")
	req.Header.Set("Authorization", "Bearer token123")
	req.Header.Set("X-Custom-Header", "custom-value")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
		t.Logf("Response: %s", w.Body.String())
	}

	var resp model.ParamTestResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(resp.Headers["User-Agent"]) > 0 && resp.Headers["User-Agent"][0] != "TestAgent/1.0" {
		t.Errorf("User-Agent = %q, want %q", resp.Headers["User-Agent"][0], "TestAgent/1.0")
	}
}

func TestParamsHeader_HTTP_CaseInsensitive(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		HeaderSummary: "Header parameters",
		Tag:           "Params",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodGet, "/params/header", nil)
	req.Header.Set("user-agent", "LowerCaseAgent/1.0")
	req.Header.Set("AUTHORIZATION", "Bearer UPPERCASE")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestParamsCookie_HTTP_ValidCookies(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		CookieSummary: "Cookie parameters",
		Tag:           "Params",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodGet, "/params/cookie", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "abc123"})
	req.AddCookie(&http.Cookie{Name: "token", Value: "xyz789"})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
		t.Logf("Response: %s", w.Body.String())
	}

	var resp model.ParamTestResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Cookies["session_id"] != "abc123" {
		t.Errorf("session_id = %q, want %q", resp.Cookies["session_id"], "abc123")
	}
}

func TestParamsCookie_HTTP_MultipleCookies(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		CookieSummary: "Cookie parameters",
		Tag:           "Params",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodGet, "/params/cookie", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "session1"})
	req.AddCookie(&http.Cookie{Name: "token", Value: "token1"})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
	}

	var resp model.ParamTestResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(resp.Cookies) != 2 {
		t.Errorf("cookies count = %d, want 2", len(resp.Cookies))
	}
}

func TestParamsMixed_HTTP_AllTypes(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		MixedSummary: "Mixed parameters",
		Tag:          "Params",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodGet, "/params/mixed/456?name=Alice", nil)
	req.Header.Set("Authorization", "Bearer mixedtoken")
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "mixed123"})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
		t.Logf("Response: %s", w.Body.String())
	}

	var resp model.ParamTestResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.PathParams["id"] != "456" {
		t.Errorf("path id = %q, want %q", resp.PathParams["id"], "456")
	}

	if len(resp.QueryParams["name"]) > 0 && resp.QueryParams["name"][0] != "Alice" {
		t.Errorf("query name = %q, want %q", resp.QueryParams["name"][0], "Alice")
	}

	if len(resp.Headers["Authorization"]) > 0 && resp.Headers["Authorization"][0] != "Bearer mixedtoken" {
		t.Errorf("Authorization = %q, want %q", resp.Headers["Authorization"][0], "Bearer mixedtoken")
	}

	if resp.Cookies["session_id"] != "mixed123" {
		t.Errorf("session_id = %q, want %q", resp.Cookies["session_id"], "mixed123")
	}
}

func TestParams_ParameterParsing(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		PathSummary:   "Path",
		QuerySummary:  "Query",
		HeaderSummary: "Header",
		CookieSummary: "Cookie",
		MixedSummary:  "Mixed",
		Tag:           "Params",
	}
	Register(api, msg)

	tests := []struct {
		name   string
		path   string
		setup  func(*http.Request)
		verify func(*testing.T, *model.ParamTestResponse)
	}{
		{
			"Path only",
			"/params/path/test123",
			func(req *http.Request) {},
			func(t *testing.T, resp *model.ParamTestResponse) {
				if resp.PathParams["id"] != "test123" {
					t.Errorf("id = %q, want %q", resp.PathParams["id"], "test123")
				}
			},
		},
		{
			"Query only",
			"/params/query?name=TestUser&age=25",
			func(req *http.Request) {},
			func(t *testing.T, resp *model.ParamTestResponse) {
				if len(resp.QueryParams["name"]) > 0 && resp.QueryParams["name"][0] != "TestUser" {
					t.Errorf("name = %q, want %q", resp.QueryParams["name"][0], "TestUser")
				}
			},
		},
		{
			"Headers only",
			"/params/header",
			func(req *http.Request) {
				req.Header.Set("User-Agent", "CustomAgent")
			},
			func(t *testing.T, resp *model.ParamTestResponse) {
				if len(resp.Headers["User-Agent"]) > 0 && resp.Headers["User-Agent"][0] != "CustomAgent" {
					t.Errorf("User-Agent = %q, want %q", resp.Headers["User-Agent"][0], "CustomAgent")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			tt.setup(req)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("status = %d, want 200", w.Code)
				return
			}

			var resp model.ParamTestResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}

			tt.verify(t, &resp)
		})
	}
}

func TestParams_AllEndpoints(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		PathSummary:   "Path",
		QuerySummary:  "Query",
		HeaderSummary: "Header",
		CookieSummary: "Cookie",
		MixedSummary:  "Mixed",
		Tag:           "Params",
	}
	Register(api, msg)

	tests := []struct {
		name string
		path string
	}{
		{"Path params", "/params/path/999"},
		{"Query params", "/params/query?name=Test&age=20"},
		{"Header params", "/params/header"},
		{"Cookie params", "/params/cookie"},
		{"Mixed params", "/params/mixed/888?name=Mixed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("status = %d, want 200", w.Code)
			}
		})
	}
}
