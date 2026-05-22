package pkg

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"runtime"
)

func ExampleRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func SpecPath(name string) string {
	return filepath.Join(ExampleRoot(), "specs", name)
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
