package repository

import (
	"errors"
	"time"

	"github.com/Ablebil/sea-catering-be/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepositoryItf interface {
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByRefreshToken(refreshToken string) (*entity.User, error)
	GetUserByID(id uuid.UUID) (*entity.User, error)
	CreateUser(user *entity.User) error
	UpdateUser(email string, user *entity.User) error
	AddRefreshToken(userId uuid.UUID, token string) error
	GetRefreshTokens(userId uuid.UUID) ([]entity.RefreshToken, error)
	RemoveRefreshToken(token string) error
	RemoveUnverifiedUsers() error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositoryItf {
	return &UserRepository{db}
}

func (r *UserRepository) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByRefreshToken(token string) (*entity.User, error) {
	var refreshToken entity.RefreshToken
	err := r.db.Preload("User").Where("token = ?", token).First(&refreshToken).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return refreshToken.User, nil
}

func (r *UserRepository) GetUserByID(id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("id = ?", id).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) CreateUser(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) UpdateUser(email string, user *entity.User) error {
	return r.db.Model(&entity.User{}).
		Where("email = ?", email).
		Updates(user).Error
}

func (r *UserRepository) AddRefreshToken(userId uuid.UUID, token string) error {
	return r.db.Create(&entity.RefreshToken{
		Token:  token,
		UserID: userId,
	}).Error
}

func (r *UserRepository) GetRefreshTokens(userId uuid.UUID) ([]entity.RefreshToken, error) {
	var refreshTokens []entity.RefreshToken
	err := r.db.Where("user_id = ?", userId).Find(&refreshTokens).Error
	return refreshTokens, err
}

func (r *UserRepository) RemoveRefreshToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&entity.RefreshToken{}).Error
}

func (r *UserRepository) RemoveUnverifiedUsers() error {
	return r.db.Where("verified = ? AND created_at < ?", false, time.Now().Add(-24*time.Hour)).Delete(&entity.User{}).Error
}
