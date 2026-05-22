package params

import (
	"context"
	"net/http"

	"example/development/pkg/model"

	"github.com/danielgtaylor/huma/v2"
)

type Messages struct {
	PathSummary   string
	QuerySummary  string
	HeaderSummary string
	CookieSummary string
	MixedSummary  string
	Tag           string
}

type Output struct {
	Body model.ParamTestResponse
}

func Register(api huma.API, msg Messages) {
	// Path parameters
	huma.Register(api, huma.Operation{
		OperationID: "params-path",
		Method:      http.MethodGet,
		Path:        "/params/path/{id}",
		Summary:     msg.PathSummary,
		Tags:        []string{msg.Tag},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" doc:"Resource ID"`
	}) (*Output, error) {
		return &Output{
			Body: model.ParamTestResponse{
				PathParams: map[string]string{"id": input.ID},
			},
		}, nil
	})

	// Query parameters
	huma.Register(api, huma.Operation{
		OperationID: "params-query",
		Method:      http.MethodGet,
		Path:        "/params/query",
		Summary:     msg.QuerySummary,
		Tags:        []string{msg.Tag},
	}, func(ctx context.Context, input *struct {
		Name string `query:"name" doc:"Name filter"`
		Age  int    `query:"age" doc:"Age filter"`
	}) (*Output, error) {
		queryParams := make(map[string][]string)
		if input.Name != "" {
			queryParams["name"] = []string{input.Name}
		}
		return &Output{
			Body: model.ParamTestResponse{QueryParams: queryParams},
		}, nil
	})

	// Header parameters
	huma.Register(api, huma.Operation{
		OperationID: "params-header",
		Method:      http.MethodGet,
		Path:        "/params/header",
		Summary:     msg.HeaderSummary,
		Tags:        []string{msg.Tag},
	}, func(ctx context.Context, input *struct {
		UserAgent     string `header:"User-Agent" doc:"User agent"`
		Authorization string `header:"Authorization" doc:"Authorization header"`
		CustomHeader  string `header:"X-Custom-Header" doc:"Custom header"`
	}) (*Output, error) {
		headers := make(map[string][]string)
		if input.UserAgent != "" {
			headers["User-Agent"] = []string{input.UserAgent}
		}
		if input.Authorization != "" {
			headers["Authorization"] = []string{input.Authorization}
		}
		if input.CustomHeader != "" {
			headers["X-Custom-Header"] = []string{input.CustomHeader}
		}
		return &Output{
			Body: model.ParamTestResponse{Headers: headers},
		}, nil
	})

	// Cookie parameters
	huma.Register(api, huma.Operation{
		OperationID: "params-cookie",
		Method:      http.MethodGet,
		Path:        "/params/cookie",
		Summary:     msg.CookieSummary,
		Tags:        []string{msg.Tag},
	}, func(ctx context.Context, input *struct {
		SessionID string `cookie:"session_id" doc:"Session ID"`
		Token     string `cookie:"token" doc:"Token"`
	}) (*Output, error) {
		cookies := make(map[string]string)
		if input.SessionID != "" {
			cookies["session_id"] = input.SessionID
		}
		if input.Token != "" {
			cookies["token"] = input.Token
		}
		return &Output{
			Body: model.ParamTestResponse{Cookies: cookies},
		}, nil
	})

	// Mixed parameters
	huma.Register(api, huma.Operation{
		OperationID: "params-mixed",
		Method:      http.MethodGet,
		Path:        "/params/mixed/{id}",
		Summary:     msg.MixedSummary,
		Tags:        []string{msg.Tag},
	}, func(ctx context.Context, input *struct {
		ID            string `path:"id" doc:"Resource ID"`
		Name          string `query:"name" doc:"Name filter"`
		Authorization string `header:"Authorization" doc:"Authorization"`
		SessionID     string `cookie:"session_id" doc:"Session ID"`
	}) (*Output, error) {
		return &Output{
			Body: model.ParamTestResponse{
				PathParams:  map[string]string{"id": input.ID},
				QueryParams: map[string][]string{"name": {input.Name}},
				Headers:     map[string][]string{"Authorization": {input.Authorization}},
				Cookies:     map[string]string{"session_id": input.SessionID},
			},
		}, nil
	})
}
