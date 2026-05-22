package body

import (
	"context"
	"net/http"

	"example/development/pkg/model"

	"github.com/danielgtaylor/huma/v2"
)

// Messages defines the interface for i18n messages
type Messages struct {
	JSONSummary      string
	XMLSummary       string
	FormSummary      string
	MultipartSummary string
	Tag              string
}

// Output represents body endpoint response
type Output struct {
	Body model.Response
}

// Register registers all request body APIs
func Register(api huma.API, msg Messages) {
	// JSON body
	huma.Register(api, huma.Operation{
		OperationID: "body-json",
		Method:      http.MethodPost,
		Path:        "/body/json",
		Summary:     msg.JSONSummary,
		Tags:        []string{msg.Tag},
	}, func(ctx context.Context, input *struct {
		Body model.JSONRequest
	}) (*Output, error) {
		return &Output{
			Body: model.SuccessResponse(input.Body),
		}, nil
	})

	// XML body
	huma.Register(api, huma.Operation{
		OperationID: "body-xml",
		Method:      http.MethodPost,
		Path:        "/body/xml",
		Summary:     msg.XMLSummary,
		Tags:        []string{msg.Tag},
	}, func(ctx context.Context, input *struct {
		Body model.XMLRequest `contentType:"application/xml"`
	}) (*Output, error) {
		return &Output{
			Body: model.SuccessResponse(input.Body),
		}, nil
	})

	// Form-urlencoded body
	huma.Register(api, huma.Operation{
		OperationID: "body-form",
		Method:      http.MethodPost,
		Path:        "/body/form",
		Summary:     msg.FormSummary,
		Tags:        []string{msg.Tag},
	}, func(ctx context.Context, input *struct {
		Body model.FormRequest `contentType:"application/x-www-form-urlencoded"`
	}) (*Output, error) {
		return &Output{
			Body: model.SuccessResponse(input.Body),
		}, nil
	})

	// Multipart form data
	huma.Register(api, huma.Operation{
		OperationID: "body-multipart",
		Method:      http.MethodPost,
		Path:        "/body/multipart",
		Summary:     msg.MultipartSummary,
		Tags:        []string{msg.Tag},
	}, func(ctx context.Context, input *struct {
		Body model.FormRequest `contentType:"multipart/form-data"`
	}) (*Output, error) {
		return &Output{
			Body: model.SuccessResponse(input.Body),
		}, nil
	})
}
