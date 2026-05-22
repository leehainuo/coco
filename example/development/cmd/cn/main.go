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
	TITLE       = "Coco API 示例"
	VERSION     = "1.0.0"
	DESCRIPTION = "使用 Gin + Huma 构建的完整 API 示例"
)

func main() {
	r := gin.Default()

	// 全局中间件
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

	// 注册所有 API 模块
	auth.Register(api, auth.Messages{
		PublicMessage: "这是一个公开端点，无需身份验证",
		BearerSuccess: "Bearer 令牌认证成功",
		APIKeySuccess: "API 密钥认证成功",
		BasicSuccess:  "基本认证成功",
		PublicSummary: "公开端点（无需认证）",
		BearerSummary: "Bearer 令牌保护的端点",
		APIKeySummary: "API 密钥保护的端点",
		BasicSummary:  "基本认证保护的端点",
		Tag:           "身份认证",
	})

	body.Register(api, body.Messages{
		JSONSummary:      "测试 JSON 请求体",
		XMLSummary:       "测试 XML 请求体",
		FormSummary:      "测试 Form-urlencoded 请求体",
		MultipartSummary: "测试 Multipart Form Data",
		Tag:              "请求体",
	})

	params.Register(api, params.Messages{
		PathSummary:   "测试路径参数",
		QuerySummary:  "测试查询参数",
		HeaderSummary: "测试 Header 参数",
		CookieSummary: "测试 Cookie 参数",
		MixedSummary:  "测试混合参数",
		Tag:           "参数",
	})

	methods.Register(api, methods.Messages{
		GetSummary:     "测试 GET 方法",
		PostSummary:    "测试 POST 方法",
		PutSummary:     "测试 PUT 方法",
		PatchSummary:   "测试 PATCH 方法",
		DeleteSummary:  "测试 DELETE 方法",
		HeadSummary:    "测试 HEAD 方法",
		OptionsSummary: "测试 OPTIONS 方法",
		Tag:            "HTTP 方法",
	})

	file.Register(api, file.Messages{
		SingleFileSummary:    "单文件上传",
		MultipleFilesSummary: "多文件上传",
		Tag:                  "文件上传",
	})

	responses.Register(api, responses.Messages{
		Response200Summary: "测试 200 OK 响应",
		Response201Summary: "测试 201 Created 响应",
		Response204Summary: "测试 204 No Content 响应",
		Response400Summary: "测试 400 Bad Request 响应",
		Response401Summary: "测试 401 Unauthorized 响应",
		Response403Summary: "测试 403 Forbidden 响应",
		Response404Summary: "测试 404 Not Found 响应",
		Response500Summary: "测试 500 Internal Server Error 响应",
		Tag:                "响应状态码",
	})

	debug.Register(api, debug.Messages{
		EchoSummary: "回显请求元信息",
		EchoDesc:    "返回所有请求头、查询参数和元信息，用于验证全局或单独请求头设置是否生效。",
		Tag:         "调试",
	})

	spec, err := api.OpenAPI().MarshalJSON()
	if err != nil {
		panic("failed to marshal OpenAPI spec, err: " + err.Error())
	}
	r.Any("/docs/*any", gin.WrapH(coco.New("",
		coco.Spec(spec),
		coco.Title(TITLE),
		coco.Lang("zh"),
	)))

	log.Println("🥥 API 文档: http://localhost:8000/docs")

	r.Run(":8000")
}
