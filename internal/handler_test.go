package internal

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/leehainuo/coco/internal/config"
)

func newTestHandler(c config.Config) http.Handler {
	return New(c)
}

func TestConfigJSON_DefaultValues(t *testing.T) {
	t.Parallel()

	h := newTestHandler(config.Config{
		Spec: config.Spec{Data: []byte(`{}`)},
	})

	req := httptest.NewRequest(http.MethodGet, "/config.json", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var m map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if m["lang"] != "en" {
		t.Errorf("lang = %v, want %q", m["lang"], "en")
	}
	if m["theme"] != "auto" {
		t.Errorf("theme = %v, want %q", m["theme"], "auto")
	}
}

func TestConfigJSON_CustomValues(t *testing.T) {
	t.Parallel()

	h := newTestHandler(config.Config{
		Spec: config.Spec{Data: []byte(`{}`)},
		UI: config.UI{
			Lang:  "zh",
			Theme: "dark",
		},
		Feature: config.Feature{
			Debug:  false,
			Export: true,
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/config.json", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	var m map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if m["lang"] != "zh" {
		t.Errorf("lang = %v, want %q", m["lang"], "zh")
	}
	if m["theme"] != "dark" {
		t.Errorf("theme = %v, want %q", m["theme"], "dark")
	}
	if m["enableDebug"] != false {
		t.Errorf("enableDebug = %v, want false", m["enableDebug"])
	}
	if m["enableExport"] != true {
		t.Errorf("enableExport = %v, want true", m["enableExport"])
	}
}

func TestConfigJSON_ContentType(t *testing.T) {
	t.Parallel()

	h := newTestHandler(config.Config{
		Spec: config.Spec{Data: []byte(`{}`)},
	})

	req := httptest.NewRequest(http.MethodGet, "/config.json", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	ct := w.Header().Get("Content-Type")
	if ct != "application/json; charset=utf-8" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/json; charset=utf-8")
	}
}

func TestSpecHandler_InlineData(t *testing.T) {
	t.Parallel()

	specData := []byte(`{"openapi":"3.0.0","info":{"title":"Test"}}`)
	h := newTestHandler(config.Config{
		Spec: config.Spec{Data: specData},
	})

	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	got := strings.TrimSpace(w.Body.String())
	if got != string(specData) {
		t.Errorf("body = %q, want %q", got, specData)
	}
}

func TestSpecHandler_SwaggerAlias(t *testing.T) {
	t.Parallel()

	specData := []byte(`{"swagger":"2.0"}`)
	h := newTestHandler(config.Config{
		Spec: config.Spec{Data: specData},
	})

	req := httptest.NewRequest(http.MethodGet, "/swagger.json", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	got := strings.TrimSpace(w.Body.String())
	if got != string(specData) {
		t.Errorf("body = %q, want %q", got, specData)
	}
}

func TestSpecHandler_ContentType(t *testing.T) {
	t.Parallel()

	h := newTestHandler(config.Config{
		Spec: config.Spec{Data: []byte(`{}`)},
	})

	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	ct := w.Header().Get("Content-Type")
	if ct != "application/json; charset=utf-8" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/json; charset=utf-8")
	}
}

func TestSpecHandler_FileNotFound(t *testing.T) {
	t.Parallel()

	h := newTestHandler(config.Config{
		Spec: config.Spec{Path: "/nonexistent/spec.json"},
	})

	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestServeIndex_RootPath(t *testing.T) {
	t.Parallel()

	h := newTestHandler(config.Config{
		Spec: config.Spec{Data: []byte(`{}`)},
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	body := w.Body.String()
	if !strings.Contains(body, `<base href="/">`) {
		t.Error("expected <base href=\"/\"> in response body")
	}
	if !strings.Contains(body, "<html") {
		t.Error("expected HTML content in response body")
	}
}

func TestServeIndex_SubPath(t *testing.T) {
	t.Parallel()

	h := newTestHandler(config.Config{
		Spec: config.Spec{Data: []byte(`{}`)},
	})

	req := httptest.NewRequest(http.MethodGet, "/docs/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	body := w.Body.String()
	if !strings.Contains(body, `<base href="/docs/">`) {
		t.Errorf("expected <base href=\"/docs/\">, got body:\n%s", body)
	}
}

func TestServeIndex_SubPathNoTrailingSlash(t *testing.T) {
	t.Parallel()

	h := newTestHandler(config.Config{
		Spec: config.Spec{Data: []byte(`{}`)},
	})

	req := httptest.NewRequest(http.MethodGet, "/docs/api", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	body := w.Body.String()
	if !strings.Contains(body, `<base href="/docs/">`) {
		t.Errorf("expected <base href=\"/docs/\">, got body:\n%s", body)
	}
}

func TestServeIndex_ContentType(t *testing.T) {
	t.Parallel()

	h := newTestHandler(config.Config{
		Spec: config.Spec{Data: []byte(`{}`)},
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	ct := w.Header().Get("Content-Type")
	if ct != "text/html; charset=utf-8" {
		t.Errorf("Content-Type = %q, want %q", ct, "text/html; charset=utf-8")
	}
}

func TestCatchAll_RedirectToSpec(t *testing.T) {
	t.Parallel()

	specData := []byte(`{"test":"redirect"}`)
	h := newTestHandler(config.Config{
		Spec: config.Spec{Data: specData},
	})

	tests := []struct {
		name string
		path string
	}{
		{"openapi nested", "/docs/openapi.json"},
		{"swagger nested", "/v1/swagger.json"},
		{"config nested", "/api/config.json"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
			}

			body := w.Body.String()
			ct := w.Header().Get("Content-Type")
			if !strings.Contains(ct, "application/json") {
				t.Errorf("Content-Type = %q, want application/json", ct)
			}
			if len(body) == 0 {
				t.Error("expected non-empty response body")
			}
		})
	}
}

func TestCatchAll_StaticAssets(t *testing.T) {
	t.Parallel()

	h := newTestHandler(config.Config{
		Spec: config.Spec{Data: []byte(`{}`)},
	})

	// Request the known CSS asset file from dist/assets/
	req := httptest.NewRequest(http.MethodGet, "/assets/style-B7PsLE-Y.css", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d for static asset", w.Code, http.StatusOK)
	}

	body, _ := io.ReadAll(w.Result().Body)
	if len(body) == 0 {
		t.Error("expected non-empty static asset response")
	}
}

func TestCatchAll_UnknownPathFallsBackToIndex(t *testing.T) {
	t.Parallel()

	h := newTestHandler(config.Config{
		Spec: config.Spec{Data: []byte(`{}`)},
	})

	req := httptest.NewRequest(http.MethodGet, "/some/unknown/path", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	body := w.Body.String()
	if !strings.Contains(body, "<html") {
		t.Error("expected HTML fallback for unknown path")
	}
}

func TestCatchAll_NestedAssetsPath(t *testing.T) {
	t.Parallel()

	h := newTestHandler(config.Config{
		Spec: config.Spec{Data: []byte(`{}`)},
	})

	// Request with a prefix before assets/
	req := httptest.NewRequest(http.MethodGet, "/docs/assets/style-B7PsLE-Y.css", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d for nested asset path", w.Code, http.StatusOK)
	}
}

func TestI18nJSON_NotConfigured(t *testing.T) {
	t.Parallel()

	h := newTestHandler(config.Config{
		Spec: config.Spec{Data: []byte(`{}`)},
	})

	req := httptest.NewRequest(http.MethodGet, "/i18n.json", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestI18nJSON_Configured(t *testing.T) {
	t.Parallel()

	payload := `{"name":"Français","messages":{"search":"Recherche"}}`
	h := newTestHandler(config.Config{
		Spec: config.Spec{Data: []byte(`{}`)},
		UI:   config.UI{I18n: config.I18n{I18nData: []byte(payload)}},
	})

	req := httptest.NewRequest(http.MethodGet, "/i18n.json", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	body, _ := io.ReadAll(w.Body)
	if string(body) != payload {
		t.Errorf("body = %q, want %q", body, payload)
	}
}

func TestConfigJSON_EnableCustomI18n(t *testing.T) {
	t.Parallel()

	h := newTestHandler(config.Config{
		Spec: config.Spec{Data: []byte(`{}`)},
		UI:   config.UI{I18n: config.I18n{I18nData: []byte(`{}`)}},
	})

	req := httptest.NewRequest(http.MethodGet, "/config.json", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	var m map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	if m["enableCustomI18n"] != true {
		t.Errorf("enableCustomI18n = %v, want true", m["enableCustomI18n"])
	}
}
