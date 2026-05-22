package internal

import (
	"context"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/leehainuo/coco/internal/config"
)

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
		switch {
		case c.URL != "":
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.URL, nil)
			if err != nil {
				sl.err = err
				return
			}

			client := &http.Client{Timeout: 10 * time.Second}
			res, err := client.Do(req)
			if err != nil {
				sl.err = err
				return
			}
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				sl.err = err
				return
			}

			sl.data, sl.err = io.ReadAll(res.Body)
		default:
			sl.data, sl.err = os.ReadFile(c.Path)
		}
	})

	sl.mu.RLock()
	defer sl.mu.RUnlock()
	return sl.data, sl.err
}
