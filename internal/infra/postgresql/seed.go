package postgresql

import (
	"log"
	"time"

	"github.com/Ablebil/sea-catering-be/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	log.Println("Starting database seeding...")

	if err := seedUsers(db); err != nil {
		return err
	}

	if err := seedMealPlans(db); err != nil {
		return err
	}

	if err := seedTestimonials(db); err != nil {
		return err
	}

	log.Println("Database seeding completed successfully")
	return nil
}

func seedUsers(db *gorm.DB) error {
	log.Println("Seeding users...")

	users := []entity.User{
		{
			ID:       parseUUID("0197b1a0-0000-7000-8000-000000000001"),
			Email:    "admin@example.com",
			Password: stringPtr("$2a$10$iavr.hje5PVy97JforQtDeVsaUkUdLDkdeXjyq0x7Al43u5SO0HPy"),
			Name:     "Admin",
			Verified: true,
			Role:     entity.RoleAdmin,
		},
		{
			ID:       parseUUID("a11ce001-e89b-12d3-a456-426614174001"),
			Email:    "alice@example.com",
			Password: stringPtr("$2a$10$iavr.hje5PVy97JforQtDeVsaUkUdLDkdeXjyq0x7Al43u5SO0HPy"),
			Name:     "Alice",
			Verified: true,
			Role:     entity.RoleUser,
		},
		{
			ID:       parseUUID("b0b00002-e89b-12d3-a456-426614174002"),
			Email:    "bob@example.com",
			Password: stringPtr("$2a$10$iavr.hje5PVy97JforQtDeVsaUkUdLDkdeXjyq0x7Al43u5SO0HPy"),
			Name:     "Bob",
			Verified: true,
			Role:     entity.RoleUser,
		},
		{
			ID:       parseUUID("ca201003-e89b-12d3-a456-426614174003"),
			Email:    "carol@example.com",
			Password: stringPtr("$2a$10$iavr.hje5PVy97JforQtDeVsaUkUdLDkdeXjyq0x7Al43u5SO0HPy"),
			Name:     "Carol",
			Verified: true,
			Role:     entity.RoleUser,
		},
		{
			ID:       parseUUID("da4e0004-e89b-12d3-a456-426614174004"),
			Email:    "dave@example.com",
			Password: stringPtr("$2a$10$iavr.hje5PVy97JforQtDeVsaUkUdLDkdeXjyq0x7Al43u5SO0HPy"),
			Name:     "Dave",
			Verified: true,
			Role:     entity.RoleUser,
		},
		{
			ID:       parseUUID("e5e00005-e89b-12d3-a456-426614174005"),
			Email:    "eve@example.com",
			Password: stringPtr("$2a$10$iavr.hje5PVy97JforQtDeVsaUkUdLDkdeXjyq0x7Al43u5SO0HPy"),
			Name:     "Eve",
			Verified: true,
			Role:     entity.RoleUser,
		},
	}

	for _, user := range users {
		var existingUser entity.User
		err := db.Where("email = ?", user.Email).First(&existingUser).Error

		if err == gorm.ErrRecordNotFound {
			now := time.Now()

			result := db.Exec(`
                INSERT INTO users (id, email, password, name, google_id, verified, role, created_at, updated_at) 
                VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
            `, user.ID, user.Email, user.Password, user.Name, nil, user.Verified, user.Role, now, now)

			if result.Error != nil {
				log.Printf("Error creating user %s: %v", user.Email, result.Error)
				return result.Error
			}
			log.Printf("Created user: %s", user.Email)
		} else if err != nil {
			log.Printf("Error checking user %s: %v", user.Email, err)
			return err
		} else {
			log.Printf("User %s already exists, skipping", user.Email)
		}
	}

	return nil
}

func seedMealPlans(db *gorm.DB) error {
	log.Println("Seeding meal plans...")

	mealPlans := []entity.MealPlan{
		{
			ID:          parseUUID("d1e70001-c1b0-4a8e-9a7c-141742660001"),
			Name:        "Diet Plan",
			Description: "Manage your weight without sacrificing flavor. Our Diet Plan is carefully crafted with calorie-controlled, nutritionally balanced meals to help you achieve your health goals while enjoying delicious food every day.",
			Price:       30000,
			PhotoURL:    "https://mifsarvarixwyfrfiiwa.supabase.co/storage/v1/object/public/media/meal-plans/diet-plan.jpg",
		},
		{
			ID:          parseUUID("d1e70002-c1b0-4a8e-9a7c-141742660002"),
			Name:        "Protein Plan",
			Description: "Fuel your fitness journey and maximize muscle growth. This power-packed plan is loaded with high-quality lean protein and complex carbs to support your workouts, accelerate recovery, and build strength.",
			Price:       40000,
			PhotoURL:    "https://mifsarvarixwyfrfiiwa.supabase.co/storage/v1/object/public/media/meal-plans/protein-plan.jpg",
		},
		{
			ID:          parseUUID("d1e70003-c1b0-4a8e-9a7c-141742660003"),
			Name:        "Royal Plan",
			Description: "Indulge in an epicurean experience with our Royal Plan. Featuring the finest, premium ingredients and gourmet recipes, this luxurious meal plan is crafted for those who desire exquisite flavors and unparalleled quality.",
			Price:       60000,
			PhotoURL:    "https://mifsarvarixwyfrfiiwa.supabase.co/storage/v1/object/public/media/meal-plans/royal-plan.jpg",
		},
		{
			ID:          parseUUID("d1e70004-c1b0-4a8e-9a7c-141742660004"),
			Name:        "Vegan Plan",
			Description: "Discover the vibrant world of innovative plant-based cuisine. Our Vegan Plan features a creative menu of globally-inspired, wholesome, and delicious meals that are packed with flavor.",
			Price:       30000,
			PhotoURL:    "https://mifsarvarixwyfrfiiwa.supabase.co/storage/v1/object/public/media/meal-plans/vegan-plan.jpg",
		},
	}

	for _, mealPlan := range mealPlans {
		var existingMealPlan entity.MealPlan
		err := db.Where("name = ?", mealPlan.Name).First(&existingMealPlan).Error

		if err == gorm.ErrRecordNotFound {
			now := time.Now()
			mealPlan.CreatedAt = &now
			mealPlan.UpdatedAt = &now

			if err := db.Create(&mealPlan).Error; err != nil {
				log.Printf("Error creating meal plan %s: %v", mealPlan.Name, err)
				return err
			}
			log.Printf("Created meal plan: %s", mealPlan.Name)
		} else if err != nil {
			log.Printf("Error checking meal plan %s: %v", mealPlan.Name, err)
			return err
		} else {
			log.Printf("Meal plan %s already exists, skipping", mealPlan.Name)
		}
	}

	return nil
}

func seedTestimonials(db *gorm.DB) error {
	log.Println("Seeding testimonials...")

	testimonials := []entity.Testimonial{
		{
			ID:       parseUUID("7e571001-d4c3-4b2a-8e9a-426614174001"),
			UserID:   parseUUID("a11ce001-e89b-12d3-a456-426614174001"),
			Name:     "Alice",
			Message:  "Great food and service!",
			Rating:   5,
			PhotoURL: "https://mifsarvarixwyfrfiiwa.supabase.co/storage/v1/object/public/media/testimonials/testi-1.jpg",
		},
		{
			ID:       parseUUID("7e571002-d4c3-4b2a-8e9a-426614174002"),
			UserID:   parseUUID("b0b00002-e89b-12d3-a456-426614174002"),
			Name:     "Bob",
			Message:  "Loved the vegan options.",
			Rating:   4,
			PhotoURL: "https://mifsarvarixwyfrfiiwa.supabase.co/storage/v1/object/public/media/testimonials/testi-2.jpg",
		},
		{
			ID:       parseUUID("7e571003-d4c3-4b2a-8e9a-426614174003"),
			UserID:   parseUUID("ca201003-e89b-12d3-a456-426614174003"),
			Name:     "Carol",
			Message:  "Affordable and tasty!",
			Rating:   5,
			PhotoURL: "https://mifsarvarixwyfrfiiwa.supabase.co/storage/v1/object/public/media/testimonials/testi-3.jpg",
		},
		{
			ID:       parseUUID("7e571004-d4c3-4b2a-8e9a-426614174004"),
			UserID:   parseUUID("da4e0004-e89b-12d3-a456-426614174004"),
			Name:     "Dave",
			Message:  "Fast delivery, will order again.",
			Rating:   4,
			PhotoURL: "https://mifsarvarixwyfrfiiwa.supabase.co/storage/v1/object/public/media/testimonials/testi-4.jpg",
		},
		{
			ID:       parseUUID("7e571005-d4c3-4b2a-8e9a-426614174005"),
			UserID:   parseUUID("e5e00005-e89b-12d3-a456-426614174005"),
			Name:     "Eve",
			Message:  "Portions are generous!",
			Rating:   5,
			PhotoURL: "https://mifsarvarixwyfrfiiwa.supabase.co/storage/v1/object/public/media/testimonials/testi-5.jpg",
		},
	}

	for _, testimonial := range testimonials {
		var user entity.User
		err := db.Where("id = ?", testimonial.UserID).First(&user).Error
		if err != nil {
			log.Printf("User with ID %s not found, skipping testimonial from %s", testimonial.UserID, testimonial.Name)
			continue
		}

		var existingTestimonial entity.Testimonial
		err = db.Where("name = ?", testimonial.Name).First(&existingTestimonial).Error

		if err == gorm.ErrRecordNotFound {
			now := time.Now()
			testimonial.CreatedAt = &now
			testimonial.UpdatedAt = &now

			if err := db.Create(&testimonial).Error; err != nil {
				log.Printf("Error creating testimonial from %s: %v", testimonial.Name, err)
				return err
			}
			log.Printf("Created testimonial from: %s", testimonial.Name)
		} else if err != nil {
			log.Printf("Error checking testimonial from %s: %v", testimonial.Name, err)
			return err
		} else {
			log.Printf("Testimonial from %s already exists, skipping", testimonial.Name)
		}
	}

	return nil
}

func parseUUID(uuidStr string) uuid.UUID {
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		panic("Invalid UUID: " + uuidStr)
	}
	return id
}

func stringPtr(s string) *string {
	return &s
}
