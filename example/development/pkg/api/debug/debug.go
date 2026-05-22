package debug

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// Messages defines the interface for i18n messages
type Messages struct {
	EchoSummary string
	EchoDesc    string
	Tag         string
}

// RequestInfo represents full request metadata for debugging
type RequestInfo struct {
	Method  string              `json:"method" doc:"HTTP method"`
	Path    string              `json:"path" doc:"Request path"`
	Headers map[string][]string `json:"headers" doc:"All request headers"`
	Query   map[string][]string `json:"query,omitempty" doc:"Query parameters"`
	Host    string              `json:"host" doc:"Request host"`
	Proto   string              `json:"proto" doc:"HTTP protocol version"`
}

// EchoOutput wraps the debug response
type EchoOutput struct {
	Body RequestInfo
}

// Register registers all debug/inspection APIs
func Register(api huma.API, msg Messages) {
	// Echo endpoint — returns all request metadata
	huma.Register(api, huma.Operation{
		OperationID: "debug-echo",
		Method:      http.MethodGet,
		Path:        "/debug/echo",
		Summary:     msg.EchoSummary,
		Description: msg.EchoDesc,
		Tags:        []string{msg.Tag},
	}, func(ctx context.Context, input *struct {
		Authorization string `header:"Authorization" required:"false" doc:"Authorization header"`
		APIKey        string `header:"X-API-Key" required:"false" doc:"API Key header"`
		CustomHeader  string `header:"X-Custom-Header" required:"false" doc:"Custom header for testing"`
		Name          string `query:"name" required:"false" doc:"Test query parameter"`
	}) (*EchoOutput, error) {
		// Reconstruct headers from what Huma parsed
		headers := make(map[string][]string)
		if input.Authorization != "" {
			headers["Authorization"] = []string{input.Authorization}
		}
		if input.APIKey != "" {
			headers["X-API-Key"] = []string{input.APIKey}
		}
		if input.CustomHeader != "" {
			headers["X-Custom-Header"] = []string{input.CustomHeader}
		}

		query := make(map[string][]string)
		if input.Name != "" {
			query["name"] = []string{input.Name}
		}

		return &EchoOutput{
			Body: RequestInfo{
				Method:  "GET",
				Path:    "/debug/echo",
				Headers: headers,
				Query:   query,
			},
		}, nil
	})

	// POST echo — test request body + headers together
	huma.Register(api, huma.Operation{
		OperationID: "debug-echo-post",
		Method:      http.MethodPost,
		Path:        "/debug/echo",
		Summary:     msg.EchoSummary + " (POST)",
		Description: msg.EchoDesc,
		Tags:        []string{msg.Tag},
	}, func(ctx context.Context, input *struct {
		Authorization string `header:"Authorization" required:"false" doc:"Authorization header"`
		APIKey        string `header:"X-API-Key" required:"false" doc:"API Key header"`
		CustomHeader  string `header:"X-Custom-Header" required:"false" doc:"Custom header for testing"`
		Body          any    `json:"body" doc:"Any request body"`
	}) (*EchoOutput, error) {
		headers := make(map[string][]string)
		if input.Authorization != "" {
			headers["Authorization"] = []string{input.Authorization}
		}
		if input.APIKey != "" {
			headers["X-API-Key"] = []string{input.APIKey}
		}
		if input.CustomHeader != "" {
			headers["X-Custom-Header"] = []string{input.CustomHeader}
		}

		return &EchoOutput{
			Body: RequestInfo{
				Method:  "POST",
				Path:    "/debug/echo",
				Headers: headers,
			},
		}, nil
	})
}
