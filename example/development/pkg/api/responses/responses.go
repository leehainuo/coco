package responses

import (
	"context"
	"net/http"

	"example/development/pkg/model"

	"github.com/danielgtaylor/huma/v2"
)

// Messages defines the interface for i18n messages
type Messages struct {
	Response200Summary string
	Response201Summary string
	Response204Summary string
	Response400Summary string
	Response401Summary string
	Response403Summary string
	Response404Summary string
	Response500Summary string
	Tag                string
}

// Register registers all response status code APIs
func Register(api huma.API, msg Messages) {
	// 200 OK
	huma.Register(api, huma.Operation{
		OperationID: "response-200",
		Method:      http.MethodGet,
		Path:        "/responses/200",
		Summary:     msg.Response200Summary,
		Tags:        []string{msg.Tag},
	}, func(ctx context.Context, input *struct{}) (*struct {
		Body model.Response
	}, error) {
		return &struct{ Body model.Response }{
			Body: model.SuccessResponse(map[string]string{"status": "success"}),
		}, nil
	})

	// 201 Created
	huma.Register(api, huma.Operation{
		OperationID: "response-201",
		Method:      http.MethodPost,
		Path:        "/responses/201",
		Summary:     msg.Response201Summary,
		Tags:        []string{msg.Tag},
	}, func(ctx context.Context, input *struct{}) (*struct {
		Status int `status:"201"`
		Body   model.Response
	}, error) {
		return &struct {
			Status int `status:"201"`
			Body   model.Response
		}{
			Status: 201,
			Body:   model.SuccessResponse(map[string]string{"status": "created"}),
		}, nil
	})

	// 204 No Content
	huma.Register(api, huma.Operation{
		OperationID: "response-204",
		Method:      http.MethodDelete,
		Path:        "/responses/204",
		Summary:     msg.Response204Summary,
		Tags:        []string{msg.Tag},
	}, func(ctx context.Context, input *struct{}) (*struct{}, error) {
		return &struct{}{}, nil
	})

	// 400 Bad Request
	huma.Register(api, huma.Operation{
		OperationID: "response-400",
		Method:      http.MethodGet,
		Path:        "/responses/400",
		Summary:     msg.Response400Summary,
		Tags:        []string{msg.Tag},
	}, func(ctx context.Context, input *struct{}) (*struct {
		Body model.ErrorResponse
	}, error) {
		return &struct{ Body model.ErrorResponse }{
			Body: model.ErrorResponseWith(400, "Bad Request", "Invalid parameter"),
		}, nil
	})

	// 401 Unauthorized
	huma.Register(api, huma.Operation{
		OperationID: "response-401",
		Method:      http.MethodGet,
		Path:        "/responses/401",
		Summary:     msg.Response401Summary,
		Tags:        []string{msg.Tag},
	}, func(ctx context.Context, input *struct{}) (*struct {
		Body model.ErrorResponse
	}, error) {
		return &struct{ Body model.ErrorResponse }{
			Body: model.ErrorResponseWith(401, "Unauthorized", "Authentication required"),
		}, nil
	})

	// 403 Forbidden
	huma.Register(api, huma.Operation{
		OperationID: "response-403",
		Method:      http.MethodGet,
		Path:        "/responses/403",
		Summary:     msg.Response403Summary,
		Tags:        []string{msg.Tag},
	}, func(ctx context.Context, input *struct{}) (*struct {
		Body model.ErrorResponse
	}, error) {
		return &struct{ Body model.ErrorResponse }{
			Body: model.ErrorResponseWith(403, "Forbidden", "Access denied"),
		}, nil
	})

	// 404 Not Found
	huma.Register(api, huma.Operation{
		OperationID: "response-404",
		Method:      http.MethodGet,
		Path:        "/responses/404",
		Summary:     msg.Response404Summary,
		Tags:        []string{msg.Tag},
	}, func(ctx context.Context, input *struct{}) (*struct {
		Body model.ErrorResponse
	}, error) {
		return &struct{ Body model.ErrorResponse }{
			Body: model.ErrorResponseWith(404, "Not Found", "Resource not found"),
		}, nil
	})

	// 500 Internal Server Error
	huma.Register(api, huma.Operation{
		OperationID: "response-500",
		Method:      http.MethodGet,
		Path:        "/responses/500",
		Summary:     msg.Response500Summary,
		Tags:        []string{msg.Tag},
	}, func(ctx context.Context, input *struct{}) (*struct {
		Body model.ErrorResponse
	}, error) {
		return &struct{ Body model.ErrorResponse }{
			Body: model.ErrorResponseWith(500, "Internal Server Error", "Something went wrong"),
		}, nil
	})
}
