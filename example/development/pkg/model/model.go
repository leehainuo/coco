package model

// Response represents a standard API response
type Response struct {
	Code int `json:"code"`
	Data any `json:"data"`
}

// SuccessResponse creates a successful response
func SuccessResponse(data any) Response {
	return Response{Code: 200, Data: data}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

// ErrorResponseWith creates an error response with the given parameters
func ErrorResponseWith(code int, err string, message string) ErrorResponse {
	return ErrorResponse{Code: code, Error: err, Message: message}
}

// JSONRequest represents a JSON request body
type JSONRequest struct {
	Name  string `json:"name" doc:"Name"`
	Email string `json:"email" doc:"Email"`
	Age   int    `json:"age" doc:"Age"`
}

// XMLRequest represents an XML request body
type XMLRequest struct {
	Name  string `json:"name" xml:"name" doc:"Name"`
	Email string `json:"email" xml:"email" doc:"Email"`
	Age   int    `json:"age" xml:"age" doc:"Age"`
}

// FormRequest represents a form request body
type FormRequest struct {
	Name  string `json:"name" form:"name" doc:"Name"`
	Email string `json:"email" form:"email" doc:"Email"`
	Age   int    `json:"age" form:"age" doc:"Age"`
}

// UpdateRequest represents a partial update request
type UpdateRequest struct {
	Name *string `json:"name,omitempty" doc:"Name"`
}

// MethodTestResponse represents a method test response
type MethodTestResponse struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Body   any    `json:"body,omitempty"`
}

// ParamTestResponse represents a parameter test response
type ParamTestResponse struct {
	PathParams  map[string]string   `json:"pathParams,omitempty"`
	QueryParams map[string][]string `json:"queryParams,omitempty"`
	Headers     map[string][]string `json:"headers,omitempty"`
	Cookies     map[string]string   `json:"cookies,omitempty"`
}

// FileUploadResponse represents a file upload response
type FileUploadResponse struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	Message  string `json:"message"`
}
