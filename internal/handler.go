package internal

import (
	"encoding/json"
	"io"
	"io/fs"
	"net/http"
	"strings"

	"github.com/leehainuo/coco/internal/assets"
	"github.com/leehainuo/coco/internal/config"
)

func New(c config.Config) http.Handler {
	distFS, err := fs.Sub(assets.DistFS, "dist")
	if err != nil {
		panic("failed to load dist assets: " + err.Error())
	}
	server := http.FileServer(http.FS(distFS))

	mux := http.NewServeMux()

	i18n := &i18nLoader{}
	loader := &specLoader{}

	mux.HandleFunc("/config.json", func(w http.ResponseWriter, r *http.Request) {
		lang := c.Lang
		if lang == "" {
			lang = "en"
		}
		theme := c.Theme
		if theme == "" {
			theme = "auto"
		}

		c := map[string]any{
			"title":            c.Title,
			"lang":             lang,
			"theme":            theme,
			"enableDebug":      c.Debug,
			"enableExport":     c.Export,
			"enableHistory":    c.History,
			"enableCustomI18n": i18n.configured(c),
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(c)
	})

	specHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		data, err := loader.loadSpec(c)
		if err != nil {
			http.Error(w, "failed to load spec: "+err.Error(), http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(data)
	}

	mux.HandleFunc("/openapi.json", specHandler)
	mux.HandleFunc("/swagger.json", specHandler)

	mux.HandleFunc("/i18n.json", func(w http.ResponseWriter, r *http.Request) {
		if !i18n.configured(c) {
			http.NotFound(w, r)
			return
		}

		data, err := i18n.loadI18n(c)
		if err != nil {
			http.Error(w, "failed to load i18n: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write(data)
	})

	serveIndex := func(w http.ResponseWriter, r *http.Request) {
		f, err := distFS.Open("index.html")
		if err != nil {
			http.Error(w, "failed to open index.html: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		data, err := io.ReadAll(f)
		if err != nil {
			http.Error(w, "failed to read index.html: "+err.Error(), http.StatusInternalServerError)
			return
		}

		html := string(data)

		basePath := "/"
		if r.URL.Path != "/" && r.URL.Path != "" {
			if strings.HasSuffix(r.URL.Path, "/") {
				basePath = r.URL.Path
			} else {
				i := strings.LastIndex(r.URL.Path, "/")
				if i > 0 {
					basePath = r.URL.Path[:i+1]
				}
			}
		}

		baseTag := `<base href="` + basePath + `">`
		if i := strings.Index(html, `<meta charset`); i != -1 {
			if end := strings.Index(html[i:], "/>"); end != -1 {
				insertPos := i + end + 2
				html = html[:insertPos] + "\n    " + baseTag + html[insertPos:]
			}
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		filename := path
		if idx := strings.LastIndex(path, "/"); idx >= 0 {
			filename = path[idx+1:]
		}

		if filename == "openapi.json" || filename == "swagger.json" || filename == "config.json" || filename == "i18n.json" {
			r2 := r.Clone(r.Context())
			r2.URL.Path = "/" + filename
			mux.ServeHTTP(w, r2)
			return
		}

		if path == "" || path == "/" || strings.HasSuffix(path, "/") {
			serveIndex(w, r)
			return
		}

		rel := strings.TrimPrefix(path, "/")
		if strings.Contains(path, "assets/") {
			if idx := strings.Index(path, "assets/"); idx >= 0 {
				rel = path[idx:]
			}
		}

		if _, err := fs.Stat(distFS, rel); err == nil {
			r2 := r.Clone(r.Context())
			r2.URL.Path = "/" + rel
			server.ServeHTTP(w, r2)
			return
		}

		serveIndex(w, r)
	})

	return mux
}
