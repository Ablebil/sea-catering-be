package usecase

import (
	"fmt"
	"mime/multipart"
	"path/filepath"

	testimonialRepository "github.com/Ablebil/sea-catering-be/internal/app/testimonial/repository"
	"github.com/Ablebil/sea-catering-be/internal/domain/dto"
	"github.com/Ablebil/sea-catering-be/internal/domain/entity"
	res "github.com/Ablebil/sea-catering-be/internal/infra/response"
	"github.com/Ablebil/sea-catering-be/internal/infra/supabase"
	"github.com/google/uuid"
)

type TestimonialUsecaseItf interface {
	GetAllTestimonials() ([]entity.Testimonial, *res.Err)
	CreateTestimonial(userID uuid.UUID, req dto.CreateTestimonialRequest, photo multipart.File, photoHeader *multipart.FileHeader) *res.Err
}

type TestimonialUsecase struct {
	TestimonialRepository testimonialRepository.TestimonialRepositoryItf
	supabase              supabase.SupabaseItf
}

func NewTestimonialUsecase(testimonialRepository testimonialRepository.TestimonialRepositoryItf, supabase supabase.SupabaseItf) TestimonialUsecaseItf {
	return &TestimonialUsecase{
		TestimonialRepository: testimonialRepository,
		supabase:              supabase,
	}
}

func (uc *TestimonialUsecase) GetAllTestimonials() ([]entity.Testimonial, *res.Err) {
	testimonials, err := uc.TestimonialRepository.GetAllTestimonials()
	if err != nil {
		return nil, res.ErrInternalServerError(res.FailedGetAllTestimonials)
	}

	return testimonials, nil
}

func (uc *TestimonialUsecase) CreateTestimonial(userID uuid.UUID, req dto.CreateTestimonialRequest, photo multipart.File, photoHeader *multipart.FileHeader) *res.Err {
	defer photo.Close()

	ext := filepath.Ext(photoHeader.Filename)
	fileName := fmt.Sprintf("testimonials/%s%s", uuid.New().String(), ext)
	mimeType := photoHeader.Header.Get("Content-Type")
	bucketName := "media"

	publicURL, err := uc.supabase.UploadFile(photo, bucketName, fileName, mimeType)
	if err != nil {
		return res.ErrInternalServerError(res.FailedUploadFile)
	}

	testimonial := &entity.Testimonial{
		UserID:   userID,
		Name:     req.Name,
		Message:  req.Message,
		Rating:   req.Rating,
		PhotoURL: publicURL,
	}

	if err := uc.TestimonialRepository.CreateTestimonial(testimonial); err != nil {
		go uc.supabase.DeleteFile(bucketName, []string{fileName})
		return res.ErrInternalServerError(res.FailedCreateTestimonial)
	}

	return nil
}
