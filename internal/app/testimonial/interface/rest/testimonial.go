package rest

import (
	conf "github.com/Ablebil/sea-catering-be/config"
	"github.com/Ablebil/sea-catering-be/internal/app/testimonial/usecase"
	"github.com/Ablebil/sea-catering-be/internal/domain/dto"
	res "github.com/Ablebil/sea-catering-be/internal/infra/response"
	"github.com/Ablebil/sea-catering-be/internal/middleware"
	"github.com/Ablebil/sea-catering-be/internal/pkg/helper"
	"github.com/Ablebil/sea-catering-be/internal/pkg/limiter"
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
	routerGroup.Post("/", middleware.Authentication, limiter.Testimonial(), testimonialHandler.CreateTestimonial)
}

// @Summary      Get All Testimonials
// @Description  Get all testimonials.
// @Tags         Testimonial
// @Produce      json
// @Success      200  {object}  res.Res{payload=[]dto.TestimonialResponse} "Get all testimonials successful"
// @Failure      500  {object}  res.Err "Internal Server Error"
// @Router       /testimonials/ [get]
func (h TestimonialHandler) GetAllTestimonials(ctx *fiber.Ctx) error {
	testimonials, err := h.TestimonialUsecase.GetAllTestimonials()
	if err != nil {
		return err
	}

	return res.OK(ctx, testimonials, res.GetAllTestimonialsSuccess)
}

// @Summary      Create Testimonial
// @Description  Create a new testimonial with photo upload. Only for authenticated users.
// @Tags         Testimonial
// @Accept       multipart/form-data
// @Produce      json
// @Param        name     formData string true  "Customer Name" example(John Doe)
// @Param        message  formData string true  "Review Message" example(The food was delicious)
// @Param        rating   formData int    true  "Rating (1-5)" example(5)
// @Param        photo    formData file   true  "Photo of the meal"
// @Success      201  {object}  res.Res "Create testimonial successful"
// @Failure      400  {object}  res.Err "Bad Request (validation error)"
// @Failure      401  {object}  res.Err "Unauthorized"
// @Failure      500  {object}  res.Err "Internal Server Error"
// @Security     ApiKeyAuth
// @Router       /testimonials/ [post]
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
