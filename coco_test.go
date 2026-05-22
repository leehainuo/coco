package coco

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/leehainuo/coco/internal/config"
)

func TestNew_ReturnsHandler(t *testing.T) {
	t.Parallel()

	h := New("", Spec([]byte(`{}`)))
	if h == nil {
		t.Fatal("New() returned nil handler")
	}
}

func TestNew_DefaultConfig(t *testing.T) {
	t.Parallel()

	var captured config.Config
	orig := New
	_ = orig

	// Verify defaults via config.json endpoint
	h := New("./testpath.json", Spec([]byte(`{}`)))
	req := httptest.NewRequest(http.MethodGet, "/config.json", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	var m map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Fatalf("failed to parse config.json: %v", err)
	}

	_ = captured

	if m["lang"] != "en" {
		t.Errorf("default lang = %v, want %q", m["lang"], "en")
	}
	if m["theme"] != "auto" {
		t.Errorf("default theme = %v, want %q", m["theme"], "auto")
	}
	if m["enableDebug"] != true {
		t.Errorf("default enableDebug = %v, want true", m["enableDebug"])
	}
	if m["enableExport"] != true {
		t.Errorf("default enableExport = %v, want true", m["enableExport"])
	}
}

func TestOption_Spec(t *testing.T) {
	t.Parallel()

	data := []byte(`{"openapi":"3.0.0"}`)
	var c config.Config
	opt := Spec(data)
	opt(&c)

	if string(c.Data) != string(data) {
		t.Errorf("Data = %q, want %q", c.Data, data)
	}
	if c.Path != "" {
		t.Errorf("Path = %q, want empty (Spec should clear Path)", c.Path)
	}
}

func TestOption_SpecURL(t *testing.T) {
	t.Parallel()

	url := "https://example.com/openapi.json"
	var c config.Config
	opt := SpecURL(url)
	opt(&c)

	if c.URL != url {
		t.Errorf("URL = %q, want %q", c.URL, url)
	}
	if c.Path != "" {
		t.Errorf("Path = %q, want empty (SpecURL should clear Path)", c.Path)
	}
}

func TestOption_Title(t *testing.T) {
	t.Parallel()

	var c config.Config
	Title("My API")(&c)
	if c.Title != "My API" {
		t.Errorf("Title = %q, want %q", c.Title, "My API")
	}
}

func TestOption_Theme(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		theme string
	}{
		{"light", "light"},
		{"dark", "dark"},
		{"auto", "auto"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var c config.Config
			Theme(tt.theme)(&c)
			if c.Theme != tt.theme {
				t.Errorf("Theme = %q, want %q", c.Theme, tt.theme)
			}
		})
	}
}

func TestOption_Lang(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		lang string
	}{
		{"english", "en"},
		{"chinese", "zh"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var c config.Config
			Lang(tt.lang)(&c)
			if c.Lang != tt.lang {
				t.Errorf("Lang = %q, want %q", c.Lang, tt.lang)
			}
		})
	}
}

func TestOption_EnableDebug(t *testing.T) {
	t.Parallel()

	var c config.Config
	EnableDebug(false)(&c)
	if c.Debug != false {
		t.Errorf("Debug = %v, want false", c.Debug)
	}

	EnableDebug(true)(&c)
	if c.Debug != true {
		t.Errorf("Debug = %v, want true", c.Debug)
	}
}

func TestOption_EnableExport(t *testing.T) {
	t.Parallel()

	var c config.Config
	EnableExport(false)(&c)
	if c.Export != false {
		t.Errorf("Export = %v, want false", c.Export)
	}

	EnableExport(true)(&c)
	if c.Export != true {
		t.Errorf("Export = %v, want true", c.Export)
	}
}

func TestOption_EnableHistory(t *testing.T) {
	t.Parallel()

	var c config.Config
	EnableHistory(false)(&c)
	if c.History != false {
		t.Errorf("History = %v, want false", c.History)
	}

	EnableHistory(true)(&c)
	if c.History != true {
		t.Errorf("History = %v, want true", c.History)
	}
}

func TestNew_MultipleOptions(t *testing.T) {
	t.Parallel()

	specData := []byte(`{"openapi":"3.0.0"}`)
	h := New("./original.json",
		Spec(specData),
		Title("Custom Title"),
		Theme("dark"),
		Lang("zh"),
		EnableDebug(false),
		EnableExport(false),
		EnableHistory(false),
	)

	if h == nil {
		t.Fatal("New() returned nil handler")
	}

	// Verify config reflects all options
	req := httptest.NewRequest(http.MethodGet, "/config.json", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	var m map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Fatalf("failed to parse config.json: %v", err)
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
	if m["enableExport"] != false {
		t.Errorf("enableExport = %v, want false", m["enableExport"])
	}
}

func TestNew_ServeConfigJSON(t *testing.T) {
	t.Parallel()

	h := New("",
		Spec([]byte(`{}`)),
		Lang("zh"),
		Theme("light"),
		EnableDebug(true),
		EnableExport(false),
	)

	req := httptest.NewRequest(http.MethodGet, "/config.json", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var m map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	if m["lang"] != "zh" {
		t.Errorf("lang = %v, want %q", m["lang"], "zh")
	}
	if m["theme"] != "light" {
		t.Errorf("theme = %v, want %q", m["theme"], "light")
	}
}

func TestNew_ServeSpecFromData(t *testing.T) {
	t.Parallel()

	specData := []byte(`{"openapi":"3.0.0","info":{"title":"Integration"}}`)
	h := New("", Spec(specData))

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

func TestNew_ServeIndex(t *testing.T) {
	t.Parallel()

	h := New("", Spec([]byte(`{}`)))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	body := w.Body.String()
	if !strings.Contains(body, "<html") {
		t.Error("expected HTML content")
	}
	if !strings.Contains(body, `<base href="/">`) {
		t.Error("expected <base href=\"/\"> tag")
	}
}

func TestNew_OptionOverride(t *testing.T) {
	t.Parallel()

	// Later options should override earlier ones
	h := New("",
		Spec([]byte(`{}`)),
		Theme("light"),
		Theme("dark"),
	)

	req := httptest.NewRequest(http.MethodGet, "/config.json", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	var m map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	if m["theme"] != "dark" {
		t.Errorf("theme = %v, want %q (last option should win)", m["theme"], "dark")
	}
}
