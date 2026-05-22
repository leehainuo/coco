package responses

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
		Response200Summary: "200 OK",
		Response201Summary: "201 Created",
		Response204Summary: "204 No Content",
		Response400Summary: "400 Bad Request",
		Response401Summary: "401 Unauthorized",
		Response403Summary: "403 Forbidden",
		Response404Summary: "404 Not Found",
		Response500Summary: "500 Internal Server Error",
		Tag:                "HTTP Responses",
	}

	if msg.Response200Summary != "200 OK" {
		t.Errorf("Response200Summary = %q, want %q", msg.Response200Summary, "200 OK")
	}

	if msg.Response201Summary != "201 Created" {
		t.Errorf("Response201Summary = %q, want %q", msg.Response201Summary, "201 Created")
	}

	if msg.Response204Summary != "204 No Content" {
		t.Errorf("Response204Summary = %q, want %q", msg.Response204Summary, "204 No Content")
	}

	if msg.Response400Summary != "400 Bad Request" {
		t.Errorf("Response400Summary = %q, want %q", msg.Response400Summary, "400 Bad Request")
	}

	if msg.Response401Summary != "401 Unauthorized" {
		t.Errorf("Response401Summary = %q, want %q", msg.Response401Summary, "401 Unauthorized")
	}

	if msg.Response403Summary != "403 Forbidden" {
		t.Errorf("Response403Summary = %q, want %q", msg.Response403Summary, "403 Forbidden")
	}

	if msg.Response404Summary != "404 Not Found" {
		t.Errorf("Response404Summary = %q, want %q", msg.Response404Summary, "404 Not Found")
	}

	if msg.Response500Summary != "500 Internal Server Error" {
		t.Errorf("Response500Summary = %q, want %q", msg.Response500Summary, "500 Internal Server Error")
	}

	if msg.Tag != "HTTP Responses" {
		t.Errorf("Tag = %q, want %q", msg.Tag, "HTTP Responses")
	}
}

func TestRegister_DoesNotPanic(t *testing.T) {
	t.Parallel()

	_, api := humatest.New(t)

	msg := Messages{
		Response200Summary: "200",
		Response201Summary: "201",
		Response204Summary: "204",
		Response400Summary: "400",
		Response401Summary: "401",
		Response403Summary: "403",
		Response404Summary: "404",
		Response500Summary: "500",
		Tag:                "Responses",
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
		Response200Summary: "200",
		Response201Summary: "201",
		Response204Summary: "204",
		Response400Summary: "400",
		Response401Summary: "401",
		Response403Summary: "403",
		Response404Summary: "404",
		Response500Summary: "500",
		Tag:                "tag",
	}

	fields := []string{
		msg.Response200Summary,
		msg.Response201Summary,
		msg.Response204Summary,
		msg.Response400Summary,
		msg.Response401Summary,
		msg.Response403Summary,
		msg.Response404Summary,
		msg.Response500Summary,
		msg.Tag,
	}

	for i, field := range fields {
		if field == "" {
			t.Errorf("field %d is empty", i)
		}
	}
}

func TestMessages_StatusCodes(t *testing.T) {
	t.Parallel()

	statusCodes := []struct {
		summary string
		code    int
	}{
		{"200 OK", 200},
		{"201 Created", 201},
		{"204 No Content", 204},
		{"400 Bad Request", 400},
		{"401 Unauthorized", 401},
		{"403 Forbidden", 403},
		{"404 Not Found", 404},
		{"500 Internal Server Error", 500},
	}

	for _, tc := range statusCodes {
		t.Run(tc.summary, func(t *testing.T) {
			t.Parallel()

			if tc.summary == "" {
				t.Error("summary is empty")
			}

			if tc.code < 100 || tc.code >= 600 {
				t.Errorf("invalid HTTP status code: %d", tc.code)
			}
		})
	}
}

// HTTP Integration Tests

func TestResponse200_HTTP(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		Response200Summary: "200 OK",
		Tag:                "Responses",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodGet, "/responses/200", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
		t.Logf("Response: %s", w.Body.String())
	}

	var resp model.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Code != 200 {
		t.Errorf("response code = %d, want 200", resp.Code)
	}
}

func TestResponse201_HTTP(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		Response201Summary: "201 Created",
		Tag:                "Responses",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodPost, "/responses/201", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusCreated)
		t.Logf("Response: %s", w.Body.String())
	}

	var resp model.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Code != 200 {
		t.Errorf("response code = %d, want 200", resp.Code)
	}
}

func TestResponse204_HTTP(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		Response204Summary: "204 No Content",
		Tag:                "Responses",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodDelete, "/responses/204", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusNoContent {
		t.Errorf("status code = %d, want 200 or 204", w.Code)
		t.Logf("Response: %s", w.Body.String())
	}
}

func TestResponse400_HTTP(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		Response400Summary: "400 Bad Request",
		Tag:                "Responses",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodGet, "/responses/400", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
		t.Logf("Response: %s", w.Body.String())
	}

	var resp model.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Code != 400 {
		t.Errorf("error code = %d, want 400", resp.Code)
	}
}

func TestResponse401_HTTP(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		Response401Summary: "401 Unauthorized",
		Tag:                "Responses",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodGet, "/responses/401", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
		t.Logf("Response: %s", w.Body.String())
	}

	var resp model.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Code != 401 {
		t.Errorf("error code = %d, want 401", resp.Code)
	}
}

func TestResponse403_HTTP(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		Response403Summary: "403 Forbidden",
		Tag:                "Responses",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodGet, "/responses/403", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
		t.Logf("Response: %s", w.Body.String())
	}

	var resp model.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Code != 403 {
		t.Errorf("error code = %d, want 403", resp.Code)
	}
}

func TestResponse404_HTTP(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		Response404Summary: "404 Not Found",
		Tag:                "Responses",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodGet, "/responses/404", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
		t.Logf("Response: %s", w.Body.String())
	}

	var resp model.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Code != 404 {
		t.Errorf("error code = %d, want 404", resp.Code)
	}
}

func TestResponse500_HTTP(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		Response500Summary: "500 Internal Server Error",
		Tag:                "Responses",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodGet, "/responses/500", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
		t.Logf("Response: %s", w.Body.String())
	}

	var resp model.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Code != 500 {
		t.Errorf("error code = %d, want 500", resp.Code)
	}
}

func TestResponses_SuccessCodes(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		Response200Summary: "200",
		Response201Summary: "201",
		Response204Summary: "204",
		Tag:                "Responses",
	}
	Register(api, msg)

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{"200 OK", http.MethodGet, "/responses/200", http.StatusOK},
		{"201 Created", http.MethodPost, "/responses/201", http.StatusCreated},
		{"204 No Content", http.MethodDelete, "/responses/204", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus && !(tt.wantStatus == http.StatusOK && w.Code == http.StatusNoContent) {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestResponses_ErrorCodes(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		Response400Summary: "400",
		Response401Summary: "401",
		Response403Summary: "403",
		Response404Summary: "404",
		Response500Summary: "500",
		Tag:                "Responses",
	}
	Register(api, msg)

	tests := []struct {
		name          string
		path          string
		wantErrorCode int
	}{
		{"400 Bad Request", "/responses/400", 400},
		{"401 Unauthorized", "/responses/401", 401},
		{"403 Forbidden", "/responses/403", 403},
		{"404 Not Found", "/responses/404", 404},
		{"500 Internal Server Error", "/responses/500", 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("HTTP status = %d, want %d", w.Code, http.StatusOK)
				return
			}

			var resp model.ErrorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}

			if resp.Code != tt.wantErrorCode {
				t.Errorf("error code = %d, want %d", resp.Code, tt.wantErrorCode)
			}
		})
	}
}
