package model

import (
	"encoding/json"
	"testing"
)

func TestSuccessResponse(t *testing.T) {
	t.Parallel()

	data := map[string]string{"key": "value", "status": "success"}
	resp := SuccessResponse(data)

	if resp.Code != 200 {
		t.Errorf("Code = %d, want 200", resp.Code)
	}

	if resp.Data == nil {
		t.Fatal("Data is nil, want non-nil")
	}

	dataMap, ok := resp.Data.(map[string]string)
	if !ok {
		t.Fatal("Data type assertion failed")
	}

	if dataMap["key"] != "value" {
		t.Errorf("Data[key] = %v, want %q", dataMap["key"], "value")
	}
}

func TestSuccessResponse_NilData(t *testing.T) {
	t.Parallel()

	resp := SuccessResponse(nil)

	if resp.Code != 200 {
		t.Errorf("Code = %d, want 200", resp.Code)
	}

	if resp.Data != nil {
		t.Errorf("Data = %v, want nil", resp.Data)
	}
}

func TestErrorResponseWith(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		code    int
		errMsg  string
		message string
	}{
		{"BadRequest", 400, "Bad Request", "Invalid parameter"},
		{"Unauthorized", 401, "Unauthorized", "Authentication required"},
		{"Forbidden", 403, "Forbidden", "Access denied"},
		{"NotFound", 404, "Not Found", "Resource not found"},
		{"InternalError", 500, "Internal Server Error", "Something went wrong"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resp := ErrorResponseWith(tt.code, tt.errMsg, tt.message)

			if resp.Code != tt.code {
				t.Errorf("Code = %d, want %d", resp.Code, tt.code)
			}

			if resp.Error != tt.errMsg {
				t.Errorf("Error = %q, want %q", resp.Error, tt.errMsg)
			}

			if resp.Message != tt.message {
				t.Errorf("Message = %q, want %q", resp.Message, tt.message)
			}
		})
	}
}

func TestResponse_JSONSerialization(t *testing.T) {
	t.Parallel()

	resp := SuccessResponse(map[string]interface{}{
		"id":   123,
		"name": "test",
	})

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal Response: %v", err)
	}

	var unmarshaled Response
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal Response: %v", err)
	}

	if unmarshaled.Code != 200 {
		t.Errorf("unmarshaled Code = %d, want 200", unmarshaled.Code)
	}

	if unmarshaled.Data == nil {
		t.Error("unmarshaled Data is nil")
	}
}

func TestErrorResponse_JSONSerialization(t *testing.T) {
	t.Parallel()

	errResp := ErrorResponseWith(404, "Not Found", "Resource not found")

	data, err := json.Marshal(errResp)
	if err != nil {
		t.Fatalf("failed to marshal ErrorResponse: %v", err)
	}

	var unmarshaled ErrorResponse
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal ErrorResponse: %v", err)
	}

	if unmarshaled.Code != 404 {
		t.Errorf("unmarshaled Code = %d, want 404", unmarshaled.Code)
	}

	if unmarshaled.Error != "Not Found" {
		t.Errorf("unmarshaled Error = %q, want %q", unmarshaled.Error, "Not Found")
	}

	if unmarshaled.Message != "Resource not found" {
		t.Errorf("unmarshaled Message = %q, want %q", unmarshaled.Message, "Resource not found")
	}
}

func TestJSONRequest_Structure(t *testing.T) {
	t.Parallel()

	req := JSONRequest{
		Name:  "Alice",
		Email: "alice@example.com",
		Age:   30,
	}

	if req.Name != "Alice" {
		t.Errorf("Name = %q, want %q", req.Name, "Alice")
	}

	if req.Email != "alice@example.com" {
		t.Errorf("Email = %q, want %q", req.Email, "alice@example.com")
	}

	if req.Age != 30 {
		t.Errorf("Age = %d, want 30", req.Age)
	}
}

func TestJSONRequest_JSONMarshaling(t *testing.T) {
	t.Parallel()

	req := JSONRequest{
		Name:  "Bob",
		Email: "bob@example.com",
		Age:   25,
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal JSONRequest: %v", err)
	}

	var unmarshaled JSONRequest
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal JSONRequest: %v", err)
	}

	if unmarshaled.Name != req.Name {
		t.Errorf("Name = %q, want %q", unmarshaled.Name, req.Name)
	}

	if unmarshaled.Email != req.Email {
		t.Errorf("Email = %q, want %q", unmarshaled.Email, req.Email)
	}

	if unmarshaled.Age != req.Age {
		t.Errorf("Age = %d, want %d", unmarshaled.Age, req.Age)
	}
}

func TestXMLRequest_Structure(t *testing.T) {
	t.Parallel()

	req := XMLRequest{
		Name:  "Charlie",
		Email: "charlie@example.com",
		Age:   35,
	}

	if req.Name != "Charlie" {
		t.Errorf("Name = %q, want %q", req.Name, "Charlie")
	}

	if req.Email != "charlie@example.com" {
		t.Errorf("Email = %q, want %q", req.Email, "charlie@example.com")
	}

	if req.Age != 35 {
		t.Errorf("Age = %d, want 35", req.Age)
	}
}

func TestFormRequest_Structure(t *testing.T) {
	t.Parallel()

	req := FormRequest{
		Name:  "Dave",
		Email: "dave@example.com",
		Age:   40,
	}

	if req.Name != "Dave" {
		t.Errorf("Name = %q, want %q", req.Name, "Dave")
	}

	if req.Email != "dave@example.com" {
		t.Errorf("Email = %q, want %q", req.Email, "dave@example.com")
	}

	if req.Age != 40 {
		t.Errorf("Age = %d, want 40", req.Age)
	}
}

func TestUpdateRequest_OmitEmpty(t *testing.T) {
	t.Parallel()

	name := "Updated Name"
	req := UpdateRequest{
		Name: &name,
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal UpdateRequest: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}

	if result["name"] != "Updated Name" {
		t.Errorf("name = %v, want %q", result["name"], "Updated Name")
	}
}

func TestUpdateRequest_NilOmitted(t *testing.T) {
	t.Parallel()

	req := UpdateRequest{
		Name: nil,
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal UpdateRequest: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}

	if _, exists := result["name"]; exists {
		t.Error("name field should be omitted when nil")
	}
}

func TestMethodTestResponse_Structure(t *testing.T) {
	t.Parallel()

	resp := MethodTestResponse{
		Method: "POST",
		Path:   "/api/test",
		Body:   map[string]string{"key": "value"},
	}

	if resp.Method != "POST" {
		t.Errorf("Method = %q, want %q", resp.Method, "POST")
	}

	if resp.Path != "/api/test" {
		t.Errorf("Path = %q, want %q", resp.Path, "/api/test")
	}

	if resp.Body == nil {
		t.Fatal("Body is nil")
	}
}

func TestParamTestResponse_Structure(t *testing.T) {
	t.Parallel()

	resp := ParamTestResponse{
		PathParams:  map[string]string{"id": "123"},
		QueryParams: map[string][]string{"name": {"test"}},
		Headers:     map[string][]string{"Authorization": {"Bearer token"}},
		Cookies:     map[string]string{"session": "abc"},
	}

	if resp.PathParams["id"] != "123" {
		t.Errorf("PathParams[id] = %q, want %q", resp.PathParams["id"], "123")
	}

	if len(resp.QueryParams["name"]) == 0 || resp.QueryParams["name"][0] != "test" {
		t.Error("QueryParams[name] not set correctly")
	}

	if len(resp.Headers["Authorization"]) == 0 || resp.Headers["Authorization"][0] != "Bearer token" {
		t.Error("Headers[Authorization] not set correctly")
	}

	if resp.Cookies["session"] != "abc" {
		t.Errorf("Cookies[session] = %q, want %q", resp.Cookies["session"], "abc")
	}
}

func TestFileUploadResponse_Structure(t *testing.T) {
	t.Parallel()

	resp := FileUploadResponse{
		Filename: "test.txt",
		Size:     1024,
		Message:  "Upload successful",
	}

	if resp.Filename != "test.txt" {
		t.Errorf("Filename = %q, want %q", resp.Filename, "test.txt")
	}

	if resp.Size != 1024 {
		t.Errorf("Size = %d, want 1024", resp.Size)
	}

	if resp.Message != "Upload successful" {
		t.Errorf("Message = %q, want %q", resp.Message, "Upload successful")
	}
}
