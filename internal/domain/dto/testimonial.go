package dto

type CreateTestimonialRequest struct {
	Name    string `json:"name" validate:"required,min=3,max=50"`
	Message string `json:"message" validate:"required,min=10,max=500"`
	Rating  int    `json:"rating" validate:"required,min=1,max=5"`
}
