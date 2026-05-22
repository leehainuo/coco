package methods

import (
	"context"
	"net/http"

	"example/development/pkg/model"

	"github.com/danielgtaylor/huma/v2"
)

type Messages struct {
	GetSummary     string
	PostSummary    string
	PutSummary     string
	PatchSummary   string
	DeleteSummary  string
	HeadSummary    string
	OptionsSummary string
	Tag            string
}

type Output struct {
	Body model.MethodTestResponse
}

func Register(api huma.API, msg Messages) {
	// GET
	huma.Register(api, huma.Operation{
		OperationID: "methods-get",
		Method:      http.MethodGet,
		Path:        "/methods/get",
		Summary:     msg.GetSummary,
		Tags:        []string{msg.Tag},
	}, func(ctx context.Context, input *struct{}) (*Output, error) {
		return &Output{
			Body: model.MethodTestResponse{Method: "GET", Path: "/methods/get"},
		}, nil
	})

	// POST
	huma.Register(api, huma.Operation{
		OperationID: "methods-post",
		Method:      http.MethodPost,
		Path:        "/methods/post",
		Summary:     msg.PostSummary,
		Tags:        []string{msg.Tag},
	}, func(ctx context.Context, input *struct {
		Body model.JSONRequest
	}) (*Output, error) {
		return &Output{
			Body: model.MethodTestResponse{Method: "POST", Path: "/methods/post", Body: input.Body},
		}, nil
	})

	// PUT
	huma.Register(api, huma.Operation{
		OperationID: "methods-put",
		Method:      http.MethodPut,
		Path:        "/methods/put",
		Summary:     msg.PutSummary,
		Tags:        []string{msg.Tag},
	}, func(ctx context.Context, input *struct {
		Body model.JSONRequest
	}) (*Output, error) {
		return &Output{
			Body: model.MethodTestResponse{Method: "PUT", Path: "/methods/put", Body: input.Body},
		}, nil
	})

	// PATCH
	huma.Register(api, huma.Operation{
		OperationID: "methods-patch",
		Method:      http.MethodPatch,
		Path:        "/methods/patch",
		Summary:     msg.PatchSummary,
		Tags:        []string{msg.Tag},
	}, func(ctx context.Context, input *struct {
		Body model.UpdateRequest
	}) (*Output, error) {
		return &Output{
			Body: model.MethodTestResponse{Method: "PATCH", Path: "/methods/patch", Body: input.Body},
		}, nil
	})

	// DELETE
	huma.Register(api, huma.Operation{
		OperationID: "methods-delete",
		Method:      http.MethodDelete,
		Path:        "/methods/delete",
		Summary:     msg.DeleteSummary,
		Tags:        []string{msg.Tag},
	}, func(ctx context.Context, input *struct{}) (*struct{}, error) {
		return &struct{}{}, nil
	})

	// HEAD
	huma.Register(api, huma.Operation{
		OperationID: "methods-head",
		Method:      http.MethodHead,
		Path:        "/methods/head",
		Summary:     msg.HeadSummary,
		Tags:        []string{msg.Tag},
	}, func(ctx context.Context, input *struct{}) (*struct{}, error) {
		return &struct{}{}, nil
	})

	// OPTIONS
	huma.Register(api, huma.Operation{
		OperationID: "methods-options",
		Method:      http.MethodOptions,
		Path:        "/methods/options",
		Summary:     msg.OptionsSummary,
		Tags:        []string{msg.Tag},
	}, func(ctx context.Context, input *struct{}) (*Output, error) {
		return &Output{
			Body: model.MethodTestResponse{
				Method: "OPTIONS",
				Path:   "/methods/options",
			},
		}, nil
	})
}
