package usecase

import (
	userRepository "github.com/Ablebil/sea-catering-be/internal/app/user/repository"
	"github.com/Ablebil/sea-catering-be/internal/domain/dto"
	res "github.com/Ablebil/sea-catering-be/internal/infra/response"
	"github.com/google/uuid"
)

type UserUsecaseItf interface {
	GetProfile(id uuid.UUID) (*dto.UserResponse, *res.Err)
	RemoveUnverifiedUsers() *res.Err
}

type UserUsecase struct {
	UserRepository userRepository.UserRepositoryItf
}

func NewUserUsecase(userRepository userRepository.UserRepositoryItf) UserUsecaseItf {
	return &UserUsecase{
		UserRepository: userRepository,
	}
}

func (uc *UserUsecase) GetProfile(id uuid.UUID) (*dto.UserResponse, *res.Err) {
	user, err := uc.UserRepository.GetUserByID(id)
	if err != nil {
		return nil, res.ErrInternalServerError(res.FailedGetUserProfile)
	}

	if user == nil {
		return nil, res.ErrNotFound(res.UserNotFound)
	}

	result := &dto.UserResponse{
		Name:  user.Name,
		Email: user.Email,
	}

	return result, nil
}

func (uc *UserUsecase) RemoveUnverifiedUsers() *res.Err {
	if err := uc.UserRepository.RemoveUnverifiedUsers(); err != nil {
		return res.ErrInternalServerError(res.FailedRemoveUnverifiedUsers)
	}

	return nil
}
