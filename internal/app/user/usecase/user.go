package usecase

import (
	userRepository "github.com/Ablebil/sea-catering-be/internal/app/user/repository"
	res "github.com/Ablebil/sea-catering-be/internal/infra/response"
)

type UserUsecaseItf interface {
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

func (uc *UserUsecase) RemoveUnverifiedUsers() *res.Err {
	if err := uc.UserRepository.RemoveUnverifiedUsers(); err != nil {
		return res.ErrInternalServerError(res.FailedRemoveUnverifiedUsers)
	}

	return nil
}
