package auth

import (
	"context"
	"net/http"

	"example/development/pkg/model"

	"github.com/danielgtaylor/huma/v2"
)

// Messages defines the interface for i18n messages
type Messages struct {
	PublicMessage string
	BearerSuccess string
	APIKeySuccess string
	BasicSuccess  string
	PublicSummary string
	BearerSummary string
	APIKeySummary string
	BasicSummary  string
	Tag           string
}

// Output represents authentication endpoint response
type Output struct {
	Body model.Response
}

// Register registers all authentication APIs
func Register(api huma.API, msg Messages) {
	// Public endpoint
	huma.Register(api, huma.Operation{
		OperationID: "auth-public",
		Method:      http.MethodGet,
		Path:        "/auth/public",
		Summary:     msg.PublicSummary,
		Tags:        []string{msg.Tag},
	}, func(ctx context.Context, input *struct{}) (*Output, error) {
		return &Output{
			Body: model.SuccessResponse(map[string]string{
				"message": msg.PublicMessage,
			}),
		}, nil
	})

	// Bearer token endpoint
	huma.Register(api, huma.Operation{
		OperationID: "auth-bearer",
		Method:      http.MethodGet,
		Path:        "/auth/bearer",
		Summary:     msg.BearerSummary,
		Tags:        []string{msg.Tag},
		Security:    []map[string][]string{{"bearerAuth": {}}},
	}, func(ctx context.Context, input *struct {
		Authorization string `header:"Authorization" doc:"Bearer token"`
	}) (*Output, error) {
		return &Output{
			Body: model.SuccessResponse(map[string]string{
				"message": msg.BearerSuccess,
				"token":   input.Authorization,
			}),
		}, nil
	})

	// API Key endpoint
	huma.Register(api, huma.Operation{
		OperationID: "auth-apikey",
		Method:      http.MethodGet,
		Path:        "/auth/api-key",
		Summary:     msg.APIKeySummary,
		Tags:        []string{msg.Tag},
		Security:    []map[string][]string{{"apiKeyAuth": {}}},
	}, func(ctx context.Context, input *struct {
		APIKey string `header:"X-API-Key" doc:"API Key"`
	}) (*Output, error) {
		return &Output{
			Body: model.SuccessResponse(map[string]string{
				"message": msg.APIKeySuccess,
				"api_key": input.APIKey,
			}),
		}, nil
	})

	// Basic auth endpoint
	huma.Register(api, huma.Operation{
		OperationID: "auth-basic",
		Method:      http.MethodGet,
		Path:        "/auth/basic",
		Summary:     msg.BasicSummary,
		Tags:        []string{msg.Tag},
		Security:    []map[string][]string{{"basicAuth": {}}},
	}, func(ctx context.Context, input *struct {
		Authorization string `header:"Authorization" doc:"Basic auth credentials"`
	}) (*Output, error) {
		return &Output{
			Body: model.SuccessResponse(map[string]string{
				"message": msg.BasicSuccess,
			}),
		}, nil
	})
}
