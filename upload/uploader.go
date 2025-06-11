package upload

import (
	"context"
	"mime/multipart"
)

type Uploader interface {
	Upload(ctx context.Context, file *multipart.FileHeader) (string, error)
}
