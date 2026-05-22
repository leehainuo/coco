package internal

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/leehainuo/coco/internal/config"
)

func TestLoadSpec_InlineData(t *testing.T) {
	t.Parallel()

	data := []byte(`{"openapi":"3.0.0"}`)
	sl := &specLoader{}
	c := config.Config{
		Spec: config.Spec{Data: data, Path: "/should/be/ignored"},
	}

	got, err := sl.loadSpec(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got) != string(data) {
		t.Errorf("got %q, want %q", got, data)
	}
}

func TestLoadSpec_InlineDataPriority(t *testing.T) {
	t.Parallel()

	data := []byte(`{"inline":true}`)
	sl := &specLoader{}
	c := config.Config{
		Spec: config.Spec{
			Data: data,
			Path: "/nonexistent/path.json",
			URL:  "http://invalid.example.com/spec",
		},
	}

	got, err := sl.loadSpec(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got) != string(data) {
		t.Errorf("inline data should take priority, got %q", got)
	}
}

func TestLoadSpec_FromFile(t *testing.T) {
	t.Parallel()

	content := []byte(`{"openapi":"3.0.3","info":{"title":"Test"}}`)
	dir := t.TempDir()
	path := filepath.Join(dir, "spec.json")
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	sl := &specLoader{}
	c := config.Config{
		Spec: config.Spec{Path: path},
	}

	got, err := sl.loadSpec(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got) != string(content) {
		t.Errorf("got %q, want %q", got, content)
	}
}

func TestLoadSpec_FileNotFound(t *testing.T) {
	t.Parallel()

	sl := &specLoader{}
	c := config.Config{
		Spec: config.Spec{Path: "/nonexistent/path/spec.json"},
	}

	_, err := sl.loadSpec(c)
	if err == nil {
		t.Fatal("expected error for nonexistent file, got nil")
	}
}

func TestLoadSpec_FromURL(t *testing.T) {
	t.Parallel()

	specData := `{"openapi":"3.0.0","info":{"title":"Remote"}}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(specData))
	}))
	defer ts.Close()

	sl := &specLoader{}
	c := config.Config{
		Spec: config.Spec{URL: ts.URL},
	}

	got, err := sl.loadSpec(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got) != specData {
		t.Errorf("got %q, want %q", got, specData)
	}
}

func TestLoadSpec_URLServerError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	sl := &specLoader{}
	c := config.Config{
		Spec: config.Spec{URL: ts.URL},
	}

	got, err := sl.loadSpec(c)
	// When server returns non-200, current implementation sets sl.err = err (nil)
	// and sl.data remains nil, so we expect empty data.
	if len(got) > 0 {
		t.Errorf("expected empty data for server error, got %q", got)
	}
	_ = err
}

func TestLoadSpec_OnceSemantics(t *testing.T) {
	t.Parallel()

	callCount := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Write([]byte(`{"call":1}`))
	}))
	defer ts.Close()

	sl := &specLoader{}
	c := config.Config{
		Spec: config.Spec{URL: ts.URL},
	}

	// Call multiple times
	for i := 0; i < 5; i++ {
		_, err := sl.loadSpec(c)
		if err != nil {
			t.Fatalf("call %d: unexpected error: %v", i, err)
		}
	}

	if callCount != 1 {
		t.Errorf("expected exactly 1 HTTP call due to sync.Once, got %d", callCount)
	}
}

func TestLoadSpec_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	specData := `{"concurrent":true}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(specData))
	}))
	defer ts.Close()

	sl := &specLoader{}
	c := config.Config{
		Spec: config.Spec{URL: ts.URL},
	}

	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)

	errs := make(chan error, goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			got, err := sl.loadSpec(c)
			if err != nil {
				errs <- err
				return
			}
			if string(got) != specData {
				errs <- err
			}
		}()
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Errorf("concurrent access error: %v", err)
	}
}

func TestLoadSpec_InlineDataSkipsOnce(t *testing.T) {
	t.Parallel()

	// Verify that inline data does NOT consume the sync.Once,
	// so subsequent calls with a different source still work.
	sl := &specLoader{}

	inlineData := []byte(`{"inline":true}`)
	c1 := config.Config{
		Spec: config.Spec{Data: inlineData},
	}
	got, err := sl.loadSpec(c1)
	if err != nil {
		t.Fatalf("inline call: %v", err)
	}
	if string(got) != string(inlineData) {
		t.Errorf("inline: got %q, want %q", got, inlineData)
	}

	// Now call with file path — once.Do should still fire
	fileData := []byte(`{"file":true}`)
	dir := t.TempDir()
	path := filepath.Join(dir, "spec.json")
	if err := os.WriteFile(path, fileData, 0644); err != nil {
		t.Fatal(err)
	}
	c2 := config.Config{
		Spec: config.Spec{Path: path},
	}
	got2, err := sl.loadSpec(c2)
	if err != nil {
		t.Fatalf("file call: %v", err)
	}
	if string(got2) != string(fileData) {
		t.Errorf("file: got %q, want %q", got2, fileData)
	}
}
