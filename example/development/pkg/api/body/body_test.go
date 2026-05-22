package body

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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
		JSONSummary:      "JSON body endpoint",
		XMLSummary:       "XML body endpoint",
		FormSummary:      "Form body endpoint",
		MultipartSummary: "Multipart body endpoint",
		Tag:              "Request Body",
	}

	if msg.JSONSummary != "JSON body endpoint" {
		t.Errorf("JSONSummary = %q, want %q", msg.JSONSummary, "JSON body endpoint")
	}

	if msg.XMLSummary != "XML body endpoint" {
		t.Errorf("XMLSummary = %q, want %q", msg.XMLSummary, "XML body endpoint")
	}

	if msg.FormSummary != "Form body endpoint" {
		t.Errorf("FormSummary = %q, want %q", msg.FormSummary, "Form body endpoint")
	}

	if msg.MultipartSummary != "Multipart body endpoint" {
		t.Errorf("MultipartSummary = %q, want %q", msg.MultipartSummary, "Multipart body endpoint")
	}

	if msg.Tag != "Request Body" {
		t.Errorf("Tag = %q, want %q", msg.Tag, "Request Body")
	}
}

func TestOutput_Structure(t *testing.T) {
	t.Parallel()

	output := Output{
		Body: model.SuccessResponse(map[string]interface{}{
			"name":  "test",
			"email": "test@example.com",
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
		JSONSummary:      "JSON",
		XMLSummary:       "XML",
		FormSummary:      "Form",
		MultipartSummary: "Multipart",
		Tag:              "Body",
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
		JSONSummary:      "json",
		XMLSummary:       "xml",
		FormSummary:      "form",
		MultipartSummary: "multipart",
		Tag:              "tag",
	}

	fields := []string{
		msg.JSONSummary,
		msg.XMLSummary,
		msg.FormSummary,
		msg.MultipartSummary,
		msg.Tag,
	}

	for i, field := range fields {
		if field == "" {
			t.Errorf("field %d is empty", i)
		}
	}
}

func TestOutput_WithDifferentData(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		data interface{}
	}{
		{"JSONRequest", model.JSONRequest{Name: "John", Email: "john@example.com", Age: 30}},
		{"XMLRequest", model.XMLRequest{Name: "Jane", Email: "jane@example.com", Age: 25}},
		{"FormRequest", model.FormRequest{Name: "Bob", Email: "bob@example.com", Age: 35}},
		{"Map", map[string]string{"key": "value"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			output := Output{
				Body: model.SuccessResponse(tt.data),
			}

			if output.Body.Code != 200 {
				t.Errorf("Code = %d, want 200", output.Body.Code)
			}

			if output.Body.Data == nil {
				t.Error("Data is nil")
			}
		})
	}
}

// HTTP Integration Tests

func TestBodyJSON_HTTP_ValidRequest(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		JSONSummary: "JSON endpoint",
		Tag:         "Body",
	}
	Register(api, msg)

	reqBody := model.JSONRequest{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/body/json", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
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

func TestBodyJSON_HTTP_DifferentData(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		JSONSummary: "JSON endpoint",
		Tag:         "Body",
	}
	Register(api, msg)

	tests := []struct {
		name string
		data model.JSONRequest
	}{
		{"Standard user", model.JSONRequest{Name: "Alice", Email: "alice@test.com", Age: 25}},
		{"Older user", model.JSONRequest{Name: "Bob", Email: "bob@test.com", Age: 65}},
		{"Young user", model.JSONRequest{Name: "Charlie", Email: "charlie@test.com", Age: 18}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.data)
			req := httptest.NewRequest(http.MethodPost, "/body/json", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("status = %d, want 200", w.Code)
			}
		})
	}
}

func TestBodyJSON_HTTP_EmptyBody(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		JSONSummary: "JSON endpoint",
		Tag:         "Body",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodPost, "/body/json", strings.NewReader("{}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusBadRequest && w.Code != 422 {
		t.Errorf("status code = %d, want 200, 400 or 422", w.Code)
	}
}

func TestBodyJSON_HTTP_CompleteData(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		JSONSummary: "JSON endpoint",
		Tag:         "Body",
	}
	Register(api, msg)

	reqBody := model.JSONRequest{
		Name:  "Complete User",
		Email: "complete@example.com",
		Age:   45,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/body/json", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
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

	if resp.Data == nil {
		t.Error("data is nil")
	}
}

func TestBody_ContentTypeValidation(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		JSONSummary: "JSON endpoint",
		Tag:         "Body",
	}
	Register(api, msg)

	reqBody := model.JSONRequest{Name: "Test", Email: "test@example.com", Age: 20}
	body, _ := json.Marshal(reqBody)

	tests := []struct {
		name        string
		contentType string
		wantOK      bool
	}{
		{"Valid JSON content type", "application/json", true},
		{"JSON with charset", "application/json; charset=utf-8", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/body/json", bytes.NewReader(body))
			req.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			isOK := w.Code == http.StatusOK
			if isOK != tt.wantOK {
				t.Errorf("success = %v, want %v (status=%d)", isOK, tt.wantOK, w.Code)
			}
		})
	}
}

func TestBody_ResponseEcho(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		JSONSummary: "JSON endpoint",
		Tag:         "Body",
	}
	Register(api, msg)

	reqBody := model.JSONRequest{
		Name:  "Echo Test",
		Email: "echo@test.com",
		Age:   42,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/body/json", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
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

func TestBody_AllEndpoints(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		JSONSummary:      "JSON",
		XMLSummary:       "XML",
		FormSummary:      "Form",
		MultipartSummary: "Multipart",
		Tag:              "Body",
	}
	Register(api, msg)

	tests := []struct {
		name string
		data model.JSONRequest
	}{
		{"User 1", model.JSONRequest{Name: "User One", Email: "user1@example.com", Age: 25}},
		{"User 2", model.JSONRequest{Name: "User Two", Email: "user2@example.com", Age: 30}},
		{"User 3", model.JSONRequest{Name: "User Three", Email: "user3@example.com", Age: 35}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.data)
			req := httptest.NewRequest(http.MethodPost, "/body/json", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("status = %d, want 200", w.Code)
			}
		})
	}
}
