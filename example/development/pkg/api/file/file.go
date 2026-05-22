package file

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"

	"example/development/pkg/model"

	"github.com/danielgtaylor/huma/v2"
)

// Messages defines the interface for i18n messages
type Messages struct {
	SingleFileSummary    string
	MultipleFilesSummary string
	Tag                  string
}

// Output represents file upload response
type Output struct {
	Body model.FileUploadResponse
}

// Register registers all file upload APIs
func Register(api huma.API, msg Messages) {
	// Single file upload
	huma.Register(api, huma.Operation{
		OperationID: "file-single",
		Method:      http.MethodPost,
		Path:        "/file/single",
		Summary:     msg.SingleFileSummary,
		Tags:        []string{msg.Tag},
	}, func(ctx context.Context, input *struct {
		Body struct {
			File *multipart.FileHeader `form:"file" doc:"File to upload"`
		}
	}) (*Output, error) {
		return &Output{
			Body: model.FileUploadResponse{
				Filename: input.Body.File.Filename,
				Size:     input.Body.File.Size,
				Message:  fmt.Sprintf("File %s uploaded successfully", input.Body.File.Filename),
			},
		}, nil
	})

	// Multiple files upload
	huma.Register(api, huma.Operation{
		OperationID: "file-multiple",
		Method:      http.MethodPost,
		Path:        "/file/multiple",
		Summary:     msg.MultipleFilesSummary,
		Tags:        []string{msg.Tag},
	}, func(ctx context.Context, input *struct {
		Body struct {
			Files []*multipart.FileHeader `form:"files" doc:"Files to upload"`
		}
	}) (*Output, error) {
		var totalSize int64
		filenames := make([]string, len(input.Body.Files))
		for i, f := range input.Body.Files {
			filenames[i] = f.Filename
			totalSize += f.Size
		}

		return &Output{
			Body: model.FileUploadResponse{
				Filename: fmt.Sprintf("%d files", len(filenames)),
				Size:     totalSize,
				Message:  fmt.Sprintf("Uploaded: %v", filenames),
			},
		}, nil
	})
}
