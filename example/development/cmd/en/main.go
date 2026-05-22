package main

import (
	"log"

	"example/development/pkg/api/auth"
	"example/development/pkg/api/body"
	"example/development/pkg/api/debug"
	"example/development/pkg/api/file"
	"example/development/pkg/api/methods"
	"example/development/pkg/api/params"
	"example/development/pkg/api/responses"
	"example/development/pkg/middleware"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"github.com/leehainuo/coco"
)

const (
	TITLE       = "Coco API Example"
	VERSION     = "1.0.0"
	DESCRIPTION = "Complete API example built with Gin + Huma"
)

func main() {
	r := gin.Default()

	// Global middleware
	r.Use(middleware.CORS())
	r.Use(middleware.RequestLogger())

	c := huma.DefaultConfig(TITLE, VERSION)
	c.Info.Description = DESCRIPTION
	c.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearerAuth": {
			Type:   "http",
			Scheme: "bearer",
		},
		"apiKeyAuth": {
			Type: "apiKey",
			Name: "X-API-Key",
			In:   "header",
		},
		"basicAuth": {
			Type:   "http",
			Scheme: "basic",
		},
	}

	api := humagin.New(r, c)

	// Register all API modules
	auth.Register(api, auth.Messages{
		PublicMessage: "This is a public endpoint, no authentication required",
		BearerSuccess: "Bearer token authentication successful",
		APIKeySuccess: "API key authentication successful",
		BasicSuccess:  "Basic authentication successful",
		PublicSummary: "Public Endpoint (No Auth Required)",
		BearerSummary: "Bearer Token Protected Endpoint",
		APIKeySummary: "API Key Protected Endpoint",
		BasicSummary:  "Basic Auth Protected Endpoint",
		Tag:           "Authentication",
	})

	body.Register(api, body.Messages{
		JSONSummary:      "Test JSON Request Body",
		XMLSummary:       "Test XML Request Body",
		FormSummary:      "Test Form-urlencoded Body",
		MultipartSummary: "Test Multipart Form Data",
		Tag:              "Request Body",
	})

	params.Register(api, params.Messages{
		PathSummary:   "Test Path Parameters",
		QuerySummary:  "Test Query Parameters",
		HeaderSummary: "Test Header Parameters",
		CookieSummary: "Test Cookie Parameters",
		MixedSummary:  "Test Mixed Parameters",
		Tag:           "Parameters",
	})

	methods.Register(api, methods.Messages{
		GetSummary:     "Test GET Method",
		PostSummary:    "Test POST Method",
		PutSummary:     "Test PUT Method",
		PatchSummary:   "Test PATCH Method",
		DeleteSummary:  "Test DELETE Method",
		HeadSummary:    "Test HEAD Method",
		OptionsSummary: "Test OPTIONS Method",
		Tag:            "HTTP Methods",
	})

	file.Register(api, file.Messages{
		SingleFileSummary:    "Upload Single File",
		MultipleFilesSummary: "Upload Multiple Files",
		Tag:                  "File Upload",
	})

	responses.Register(api, responses.Messages{
		Response200Summary: "Test 200 OK Response",
		Response201Summary: "Test 201 Created Response",
		Response204Summary: "Test 204 No Content Response",
		Response400Summary: "Test 400 Bad Request Response",
		Response401Summary: "Test 401 Unauthorized Response",
		Response403Summary: "Test 403 Forbidden Response",
		Response404Summary: "Test 404 Not Found Response",
		Response500Summary: "Test 500 Internal Server Error Response",
		Tag:                "Response Status Codes",
	})

	debug.Register(api, debug.Messages{
		EchoSummary: "Echo Request Metadata",
		EchoDesc:    "Returns all request headers, query params, and metadata. Use this to verify that global or per-request headers are being sent correctly.",
		Tag:         "Debug",
	})

	spec, err := api.OpenAPI().MarshalJSON()
	if err != nil {
		panic("failed to marshal OpenAPI spec, err: " + err.Error())
	}
	r.Any("/docs/*any", gin.WrapH(coco.New("",
		coco.Spec(spec),
		coco.Title(TITLE),
		coco.Lang("en"),
	)))

	log.Println("🥥 API Docs: http://localhost:8000/docs")

	r.Run(":8000")
}
