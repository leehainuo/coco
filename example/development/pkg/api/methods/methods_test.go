package methods

import (
	"bytes"
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
		GetSummary:     "GET method",
		PostSummary:    "POST method",
		PutSummary:     "PUT method",
		PatchSummary:   "PATCH method",
		DeleteSummary:  "DELETE method",
		HeadSummary:    "HEAD method",
		OptionsSummary: "OPTIONS method",
		Tag:            "HTTP Methods",
	}

	if msg.GetSummary != "GET method" {
		t.Errorf("GetSummary = %q, want %q", msg.GetSummary, "GET method")
	}

	if msg.PostSummary != "POST method" {
		t.Errorf("PostSummary = %q, want %q", msg.PostSummary, "POST method")
	}

	if msg.PutSummary != "PUT method" {
		t.Errorf("PutSummary = %q, want %q", msg.PutSummary, "PUT method")
	}

	if msg.PatchSummary != "PATCH method" {
		t.Errorf("PatchSummary = %q, want %q", msg.PatchSummary, "PATCH method")
	}

	if msg.DeleteSummary != "DELETE method" {
		t.Errorf("DeleteSummary = %q, want %q", msg.DeleteSummary, "DELETE method")
	}

	if msg.HeadSummary != "HEAD method" {
		t.Errorf("HeadSummary = %q, want %q", msg.HeadSummary, "HEAD method")
	}

	if msg.OptionsSummary != "OPTIONS method" {
		t.Errorf("OptionsSummary = %q, want %q", msg.OptionsSummary, "OPTIONS method")
	}

	if msg.Tag != "HTTP Methods" {
		t.Errorf("Tag = %q, want %q", msg.Tag, "HTTP Methods")
	}
}

func TestOutput_Structure(t *testing.T) {
	t.Parallel()

	output := Output{
		Body: model.MethodTestResponse{
			Method: "GET",
			Path:   "/methods/get",
			Body:   nil,
		},
	}

	if output.Body.Method != "GET" {
		t.Errorf("Method = %q, want %q", output.Body.Method, "GET")
	}

	if output.Body.Path != "/methods/get" {
		t.Errorf("Path = %q, want %q", output.Body.Path, "/methods/get")
	}
}

func TestRegister_DoesNotPanic(t *testing.T) {
	t.Parallel()

	_, api := humatest.New(t)

	msg := Messages{
		GetSummary:     "GET",
		PostSummary:    "POST",
		PutSummary:     "PUT",
		PatchSummary:   "PATCH",
		DeleteSummary:  "DELETE",
		HeadSummary:    "HEAD",
		OptionsSummary: "OPTIONS",
		Tag:            "Methods",
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
		GetSummary:     "get",
		PostSummary:    "post",
		PutSummary:     "put",
		PatchSummary:   "patch",
		DeleteSummary:  "delete",
		HeadSummary:    "head",
		OptionsSummary: "options",
		Tag:            "tag",
	}

	fields := []string{
		msg.GetSummary,
		msg.PostSummary,
		msg.PutSummary,
		msg.PatchSummary,
		msg.DeleteSummary,
		msg.HeadSummary,
		msg.OptionsSummary,
		msg.Tag,
	}

	for i, field := range fields {
		if field == "" {
			t.Errorf("field %d is empty", i)
		}
	}
}

func TestOutput_WithBody(t *testing.T) {
	t.Parallel()

	requestBody := model.JSONRequest{
		Name:  "Test User",
		Email: "test@example.com",
		Age:   25,
	}

	output := Output{
		Body: model.MethodTestResponse{
			Method: "POST",
			Path:   "/methods/post",
			Body:   requestBody,
		},
	}

	if output.Body.Method != "POST" {
		t.Errorf("Method = %q, want %q", output.Body.Method, "POST")
	}

	if output.Body.Body == nil {
		t.Error("Body is nil")
	}
}

func TestOutput_AllMethods(t *testing.T) {
	t.Parallel()

	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			t.Parallel()

			output := Output{
				Body: model.MethodTestResponse{
					Method: method,
					Path:   "/methods/" + method,
				},
			}

			if output.Body.Method != method {
				t.Errorf("Method = %q, want %q", output.Body.Method, method)
			}
		})
	}
}

// HTTP Integration Tests

func TestMethodGET_HTTP(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		GetSummary: "GET method",
		Tag:        "Methods",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodGet, "/methods/get", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
		t.Logf("Response: %s", w.Body.String())
	}

	var resp model.MethodTestResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Method != "GET" {
		t.Errorf("method = %q, want %q", resp.Method, "GET")
	}

	if resp.Path != "/methods/get" {
		t.Errorf("path = %q, want %q", resp.Path, "/methods/get")
	}
}

func TestMethodPOST_HTTP_WithBody(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		PostSummary: "POST method",
		Tag:         "Methods",
	}
	Register(api, msg)

	reqBody := model.JSONRequest{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/methods/post", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
		t.Logf("Response: %s", w.Body.String())
	}

	var resp model.MethodTestResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Method != "POST" {
		t.Errorf("method = %q, want %q", resp.Method, "POST")
	}
}

func TestMethodPOST_HTTP_EmptyBody(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		PostSummary: "POST method",
		Tag:         "Methods",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodPost, "/methods/post", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != 422 {
		t.Errorf("status code = %d, want 200 or 422", w.Code)
	}
}

func TestMethodPUT_HTTP(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		PutSummary: "PUT method",
		Tag:        "Methods",
	}
	Register(api, msg)

	reqBody := model.JSONRequest{
		Name:  "Updated User",
		Email: "updated@example.com",
		Age:   35,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/methods/put", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
		t.Logf("Response: %s", w.Body.String())
	}

	var resp model.MethodTestResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Method != "PUT" {
		t.Errorf("method = %q, want %q", resp.Method, "PUT")
	}
}

func TestMethodPATCH_HTTP(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		PatchSummary: "PATCH method",
		Tag:          "Methods",
	}
	Register(api, msg)

	reqBody := model.UpdateRequest{
		Name: strPtr("Patched User"),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPatch, "/methods/patch", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
		t.Logf("Response: %s", w.Body.String())
	}

	var resp model.MethodTestResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Method != "PATCH" {
		t.Errorf("method = %q, want %q", resp.Method, "PATCH")
	}
}

func TestMethodDELETE_HTTP(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		DeleteSummary: "DELETE method",
		Tag:           "Methods",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodDelete, "/methods/delete", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusNoContent {
		t.Errorf("status code = %d, want 200 or 204", w.Code)
		t.Logf("Response: %s", w.Body.String())
	}
}

func TestMethodHEAD_HTTP(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		HeadSummary: "HEAD method",
		Tag:         "Methods",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodHead, "/methods/head", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
	}

	if w.Body.Len() != 0 {
		t.Errorf("HEAD response should have no body, got %d bytes", w.Body.Len())
	}
}

func TestMethodOPTIONS_HTTP(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		OptionsSummary: "OPTIONS method",
		Tag:            "Methods",
	}
	Register(api, msg)

	req := httptest.NewRequest(http.MethodOptions, "/methods/options", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusNoContent {
		t.Errorf("status code = %d, want 200 or 204", w.Code)
		t.Logf("Response: %s", w.Body.String())
	}

	if w.Body.Len() > 0 {
		var resp model.MethodTestResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.Method != "OPTIONS" {
			t.Errorf("method = %q, want %q", resp.Method, "OPTIONS")
		}
	}
}

func TestMethods_AllInSequence(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		GetSummary:     "GET",
		PostSummary:    "POST",
		PutSummary:     "PUT",
		PatchSummary:   "PATCH",
		DeleteSummary:  "DELETE",
		HeadSummary:    "HEAD",
		OptionsSummary: "OPTIONS",
		Tag:            "Methods",
	}
	Register(api, msg)

	tests := []struct {
		method string
		path   string
		body   []byte
	}{
		{http.MethodGet, "/methods/get", nil},
		{http.MethodPost, "/methods/post", []byte(`{"name":"Test","email":"test@example.com","age":25}`)},
		{http.MethodPut, "/methods/put", []byte(`{"name":"Test","email":"test@example.com","age":25}`)},
		{http.MethodPatch, "/methods/patch", []byte(`{"name":"Test"}`)},
		{http.MethodDelete, "/methods/delete", nil},
		{http.MethodHead, "/methods/head", nil},
		{http.MethodOptions, "/methods/options", nil},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			var req *http.Request
			if tt.body != nil {
				req = httptest.NewRequest(tt.method, tt.path, bytes.NewReader(tt.body))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tt.method, tt.path, nil)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK && w.Code != http.StatusNoContent {
				t.Errorf("status = %d, want 200 or 204", w.Code)
			}
		})
	}
}

func TestMethods_ResponseBodyValidation(t *testing.T) {
	t.Parallel()

	r, api := setupTestServer(t)

	msg := Messages{
		PostSummary: "POST method",
		Tag:         "Methods",
	}
	Register(api, msg)

	reqBody := model.JSONRequest{
		Name:  "Validation Test",
		Email: "validation@example.com",
		Age:   28,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/methods/post", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", w.Code, http.StatusOK)
	}

	var resp model.MethodTestResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Method != "POST" {
		t.Errorf("method = %q, want %q", resp.Method, "POST")
	}

	if resp.Path != "/methods/post" {
		t.Errorf("path = %q, want %q", resp.Path, "/methods/post")
	}

	if resp.Body == nil {
		t.Error("body is nil")
	}
}

func strPtr(s string) *string {
	return &s
}
