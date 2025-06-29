package helper

import (
	"mime/multipart"
	"time"

	res "github.com/Ablebil/sea-catering-be/internal/infra/response"
)

type HelperItf interface {
	ValidateImageFile(file multipart.File, header *multipart.FileHeader, maxSize int64) *res.Err
	ParseDateRange(startDate string, endDate string) (time.Time, time.Time, *res.Err)
}

type Helper struct{}

func NewHelper() HelperItf {
	return &Helper{}
}
