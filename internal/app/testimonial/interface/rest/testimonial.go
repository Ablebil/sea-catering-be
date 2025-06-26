package rest

import (
	conf "github.com/Ablebil/sea-catering-be/config"
	"github.com/Ablebil/sea-catering-be/internal/app/testimonial/usecase"
	"github.com/Ablebil/sea-catering-be/internal/domain/dto"
	res "github.com/Ablebil/sea-catering-be/internal/infra/response"
	"github.com/Ablebil/sea-catering-be/internal/middleware"
	"github.com/Ablebil/sea-catering-be/internal/pkg/helper"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TestimonialHandler struct {
	Validator          *validator.Validate
	TestimonialUsecase usecase.TestimonialUsecaseItf
	helper             helper.HelperItf
	conf               *conf.Config
}

func NewTestimonialHandler(routerGroup fiber.Router, validator *validator.Validate, testimonialUsecase usecase.TestimonialUsecaseItf, middleware middleware.MiddlewareItf, helper helper.HelperItf, conf *conf.Config) {
	testimonialHandler := TestimonialHandler{
		Validator:          validator,
		TestimonialUsecase: testimonialUsecase,
		helper:             helper,
		conf:               conf,
	}

	routerGroup = routerGroup.Group("/testimonials")
	routerGroup.Get("/", testimonialHandler.GetAllTestimonials)
	routerGroup.Post("/", middleware.Authentication, testimonialHandler.CreateTestimonial)
}

func (h TestimonialHandler) GetAllTestimonials(ctx *fiber.Ctx) error {
	testimonials, err := h.TestimonialUsecase.GetAllTestimonials()
	if err != nil {
		return err
	}

	return res.OK(ctx, testimonials, res.GetAllTestimonialsSuccess)
}

func (h TestimonialHandler) CreateTestimonial(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uuid.UUID)
	req := new(dto.CreateTestimonialRequest)

	if err := ctx.BodyParser(req); err != nil {
		return res.ErrBadRequest(res.InvalidFormData)
	}

	if err := h.Validator.Struct(req); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return res.ErrInternalServerError(res.FailedValidateRequest)
		}

		return res.ErrValidation(validationErrors)
	}

	fileHeader, err := ctx.FormFile("photo")
	if err != nil {
		return res.ErrBadRequest(res.FileIsRequired)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return res.ErrInternalServerError(res.FailedToOpenFile)
	}

	maxSize := int64(h.conf.MaxFileSize) * 1024 * 1024
	if err := h.helper.ValidateImageFile(file, fileHeader, maxSize); err != nil {
		file.Close()
		return err
	}

	if err := h.TestimonialUsecase.CreateTestimonial(userID, *req, file, fileHeader); err != nil {
		return err
	}

	return res.Created(ctx, nil, res.CreateTestimonialSuccess)
}
