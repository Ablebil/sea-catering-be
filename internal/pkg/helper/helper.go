package helper

import (
	"mime/multipart"

	res "github.com/Ablebil/sea-catering-be/internal/infra/response"
)

type HelperItf interface {
	ValidateImageFile(file multipart.File, header *multipart.FileHeader, maxSize int64) *res.Err
}

type Helper struct{}

func NewHelper() HelperItf {
	return &Helper{}
}

func (h *Helper) ValidateImageFile(file multipart.File, header *multipart.FileHeader, maxSize int64) *res.Err {
	return ValidateImageFile(file, header, maxSize)
}
