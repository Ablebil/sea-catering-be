package repository

import (
	"github.com/Ablebil/sea-catering-be/internal/domain/entity"
	"gorm.io/gorm"
)

type TestimonialRepositoryItf interface {
	GetAllTestimonials() ([]entity.Testimonial, error)
	CreateTestimonial(testimonial *entity.Testimonial) error
}

type TestimonialRepository struct {
	db *gorm.DB
}

func NewTestimonialRepository(db *gorm.DB) TestimonialRepositoryItf {
	return &TestimonialRepository{
		db: db,
	}
}

func (r *TestimonialRepository) GetAllTestimonials() ([]entity.Testimonial, error) {
	var testimonials []entity.Testimonial
	err := r.db.Order("created_at desc").Find(&testimonials).Error
	return testimonials, err
}

func (r *TestimonialRepository) CreateTestimonial(testimonial *entity.Testimonial) error {
	return r.db.Create(testimonial).Error
}
