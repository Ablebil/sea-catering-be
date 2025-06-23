package usecase

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	conf "github.com/Ablebil/sea-catering-be/config"
	userRepository "github.com/Ablebil/sea-catering-be/internal/app/user/repository"
	"github.com/Ablebil/sea-catering-be/internal/domain/dto"
	"github.com/Ablebil/sea-catering-be/internal/domain/entity"
	"github.com/Ablebil/sea-catering-be/internal/infra/email"
	"github.com/Ablebil/sea-catering-be/internal/infra/jwt"
	"github.com/Ablebil/sea-catering-be/internal/infra/redis"
	res "github.com/Ablebil/sea-catering-be/internal/infra/response"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUsecaseItf interface {
	Register(req dto.RegisterRequest) *res.Err
	VerifyOTP(req dto.VerifyOTPRequest) (string, string, *res.Err)
	Login(req dto.LoginRequest) (string, string, *res.Err)
	RefreshToken(req dto.RefreshTokenRequest) (string, string, *res.Err)
}

type AuthUsecase struct {
	userRepository userRepository.UserRepositoryItf
	db             *gorm.DB
	conf           *conf.Config
	jwt            jwt.JWTItf
	email          email.EmailItf
	redis          redis.RedisItf
}

func NewAuthUsecase(userRepository userRepository.UserRepositoryItf, db *gorm.DB, conf *conf.Config, jwt jwt.JWTItf, email email.EmailItf, redis redis.RedisItf) AuthUsecaseItf {
	return &AuthUsecase{
		userRepository: userRepository,
		db:             db,
		conf:           conf,
		jwt:            jwt,
		email:          email,
		redis:          redis,
	}
}

func (uc *AuthUsecase) Register(req dto.RegisterRequest) *res.Err {
	user, err := uc.userRepository.GetUserByEmail(req.Email)
	if err != nil {
		return res.ErrInternalServerError(res.FailedFindUser)
	}

	if user != nil {
		return res.ErrConflict(res.EmailAlreadyExists)
	} else {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return res.ErrInternalServerError(res.FailedHashPassword)
		}

		hashedPassword := string(hashed)

		newUser := &entity.User{
			Name:     req.Name,
			Email:    req.Email,
			Password: &hashedPassword,
		}

		if err := uc.userRepository.CreateUser(newUser); err != nil {
			return res.ErrInternalServerError(res.FailedCreateUser)
		}
	}

	max := big.NewInt(900000)
	n, randErr := rand.Int(rand.Reader, max)
	if randErr != nil {
		return res.ErrInternalServerError(res.FailedGenerateOTP)
	}

	otp := fmt.Sprintf("%06d", 100000+n.Int64())
	expiration := 5 * time.Minute

	if err := uc.redis.SetOTP(req.Email, otp, expiration); err != nil {
		return res.ErrInternalServerError(res.FailedStoreOTP)
	}

	if err := uc.email.SendOTPEmail(req.Email, otp); err != nil {
		return res.ErrInternalServerError(res.FailedSendOTPEmail)
	}

	return nil
}

func (uc *AuthUsecase) VerifyOTP(req dto.VerifyOTPRequest) (string, string, *res.Err) {
	user, err := uc.userRepository.GetUserByEmail(req.Email)
	if err != nil {
		return "", "", res.ErrInternalServerError(res.FailedFindUser)
	}

	if user == nil {
		return "", "", res.ErrNotFound(res.UserNotFound)
	}

	storedOTP, err := uc.redis.GetOTP(req.Email)
	if err != nil || storedOTP != req.OTP {
		return "", "", res.ErrBadRequest(res.InvalidOTP)
	}

	if err := uc.redis.DeleteOTP(req.Email); err != nil {
		return "", "", res.ErrInternalServerError(res.FailedDeleteOTP)
	}

	refreshToken, err := uc.jwt.GenerateRefershToken(user.ID, false)
	if err != nil {
		return "", "", res.ErrInternalServerError(res.FailedGenerateRefreshToken)
	}

	if err := uc.userRepository.AddRefreshToken(user.ID, refreshToken); err != nil {
		return "", "", res.ErrInternalServerError(res.FailedAddRefreshToken)
	}

	user.Verified = true

	if err := uc.userRepository.UpdateUser(req.Email, user); err != nil {
		return "", "", res.ErrInternalServerError(res.FailedUpdateUser)
	}

	accessToken, err := uc.jwt.GenerateAccessToken(user.ID, user.Name, user.Email)
	if err != nil {
		return "", "", res.ErrInternalServerError(res.FailedGenerateAccessToken)
	}

	return accessToken, refreshToken, nil
}

func (uc *AuthUsecase) Login(req dto.LoginRequest) (string, string, *res.Err) {
	user, err := uc.userRepository.GetUserByEmail(req.Email)
	if err != nil {
		return "", "", res.ErrInternalServerError(res.FailedFindUser)
	}

	if user == nil {
		return "", "", res.ErrUnauthorized(res.InvalidCredentials)
	}

	if bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(req.Password)) != nil {
		return "", "", res.ErrUnauthorized(res.InvalidCredentials)
	}

	if !user.Verified {
		return "", "", res.ErrUnauthorized(res.UserNotVerified)
	}

	refreshToken, err := uc.jwt.GenerateRefershToken(user.ID, false)
	if err != nil {
		return "", "", res.ErrInternalServerError(res.FailedGenerateRefreshToken)
	}

	refreshTokens, err := uc.userRepository.GetRefreshTokens(user.ID)
	if err != nil {
		return "", "", res.ErrInternalServerError(res.FailedGetRefreshTokens)
	}

	if len(refreshTokens) >= 2 {
		if err := uc.userRepository.RemoveRefreshToken(refreshTokens[0].Token); err != nil {
			return "", "", res.ErrInternalServerError(res.FailedRemoveRefreshToken)
		}
	}

	if err := uc.userRepository.AddRefreshToken(user.ID, refreshToken); err != nil {
		return "", "", res.ErrInternalServerError(res.FailedAddRefreshToken)
	}

	accessToken, err := uc.jwt.GenerateAccessToken(user.ID, user.Name, user.Email)
	if err != nil {
		return "", "", res.ErrInternalServerError(res.FailedGenerateAccessToken)
	}

	return accessToken, refreshToken, nil
}

func (uc *AuthUsecase) RefreshToken(req dto.RefreshTokenRequest) (string, string, *res.Err) {
	user, err := uc.userRepository.GetUserByRefreshToken(req.RefreshToken)
	if err != nil {
		return "", "", res.ErrInternalServerError(res.FailedFindUser)
	}

	if user == nil {
		return "", "", res.ErrUnauthorized(res.InvalidRefreshToken)
	}

	if _, err := uc.jwt.VerifyRefreshToken(req.RefreshToken); err != nil {
		return "", "", res.ErrUnauthorized(res.InvalidRefreshToken)
	}

	accessToken, err := uc.jwt.GenerateAccessToken(user.ID, user.Name, user.Email)
	if err != nil {
		return "", "", res.ErrInternalServerError(res.FailedGenerateAccessToken)
	}

	refreshToken, err := uc.jwt.GenerateRefershToken(user.ID, false)
	if err != nil {
		return "", "", res.ErrInternalServerError(res.FailedGenerateRefreshToken)
	}

	if err := uc.userRepository.RemoveRefreshToken(req.RefreshToken); err != nil {
		return "", "", res.ErrInternalServerError(res.FailedRemoveRefreshToken)
	}

	if err := uc.userRepository.AddRefreshToken(user.ID, refreshToken); err != nil {
		return "", "", res.ErrInternalServerError(res.FailedAddRefreshToken)
	}

	return accessToken, refreshToken, nil
}
