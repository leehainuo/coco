package internal

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/leehainuo/coco/internal/config"
)

// source describes where a resource can be loaded from. Resolution
// priority is: inline data, then URL, then file path.
type source struct {
	data []byte
	url  string
	path string
}

// loadSource fetches the resource described by src. Inline data is
// returned as-is; a URL is fetched over HTTP; otherwise the file at
// path is read.
func loadSource(src source) ([]byte, error) {
	switch {
	case len(src.data) > 0:
		return src.data, nil
	case src.url != "":
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, src.url, nil)
		if err != nil {
			return nil, err
		}

		client := &http.Client{Timeout: 10 * time.Second}
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return nil, errors.New("unexpected status: " + res.Status)
		}

		return io.ReadAll(res.Body)
	case src.path != "":
		return os.ReadFile(src.path)
	default:
		return nil, nil
	}
}

type specLoader struct {
	mu   sync.RWMutex
	once sync.Once
	data []byte
	err  error
}

// loadSpec - loads the OpenAPI spec from the configured source.
// Priority:
// 1. Inline data (c.Data)
// 2. URL (c.URL)
// 3. File path (c.Path)
func (sl *specLoader) loadSpec(c config.Config) ([]byte, error) {
	if len(c.Data) > 0 {
		return c.Data, nil
	}

	sl.once.Do(func() {
		sl.data, sl.err = loadSource(source{url: c.URL, path: c.Path})
	})

	sl.mu.RLock()
	defer sl.mu.RUnlock()
	return sl.data, sl.err
}

type i18nLoader struct {
	mu   sync.RWMutex
	once sync.Once
	data []byte
	err  error
}

// configured reports whether a custom i18n source has been provided.
func (il *i18nLoader) configured(c config.Config) bool {
	return len(c.I18nData) > 0 || c.I18nURL != "" || c.I18nPath != ""
}

// loadI18n - loads the custom translation file from the configured
// source, mirroring loadSpec's resolution priority.
func (il *i18nLoader) loadI18n(c config.Config) ([]byte, error) {
	if len(c.I18nData) > 0 {
		return c.I18nData, nil
	}

	il.once.Do(func() {
		il.data, il.err = loadSource(source{url: c.I18nURL, path: c.I18nPath})
	})

	il.mu.RLock()
	defer il.mu.RUnlock()
	return il.data, il.err
}
