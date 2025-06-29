package helper

import (
	"io"
	"mime/multipart"
	"net/http"

	res "github.com/Ablebil/sea-catering-be/internal/infra/response"
)

func (h *Helper) ValidateImageFile(file multipart.File, header *multipart.FileHeader, maxSize int64) *res.Err {
	if header.Size > maxSize {
		return res.ErrEntityTooLarge(res.FileSizeExceedsLimit)
	}

	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return res.ErrInternalServerError(res.FailedReadFileForValidation)
	}

	contentType := http.DetectContentType(buffer)

	allowedMimeTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
	}

	if !allowedMimeTypes[contentType] {
		return res.ErrUnprocessableEntity(res.InvalidFileType)
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return res.ErrInternalServerError(res.FailedToResetFileReadPointer)
	}

	return nil
}
