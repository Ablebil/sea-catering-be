package usecase

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"

	conf "github.com/Ablebil/sea-catering-be/config"
	userRepository "github.com/Ablebil/sea-catering-be/internal/app/user/repository"
	"github.com/Ablebil/sea-catering-be/internal/domain/dto"
	"github.com/Ablebil/sea-catering-be/internal/domain/entity"
	"github.com/Ablebil/sea-catering-be/internal/infra/email"
	"github.com/Ablebil/sea-catering-be/internal/infra/jwt"
	"github.com/Ablebil/sea-catering-be/internal/infra/oauth"
	"github.com/Ablebil/sea-catering-be/internal/infra/redis"
	res "github.com/Ablebil/sea-catering-be/internal/infra/response"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUsecaseItf interface {
	Register(req dto.RegisterRequest) *res.Err
	VerifyOTP(req dto.VerifyOTPRequest) (string, string, *res.Err)
	Login(req dto.LoginRequest) (string, string, *res.Err)
	GoogleLogin() (string, *res.Err)
	GoogleCallback(req *dto.GoogleCallbackRequest) (string, string, bool, *res.Err)
	RefreshToken(req dto.RefreshTokenRequest) (string, string, *res.Err)
	Logout(req dto.LogoutRequest) *res.Err
}

type AuthUsecase struct {
	userRepository userRepository.UserRepositoryItf
	db             *gorm.DB
	conf           *conf.Config
	jwt            jwt.JWTItf
	email          email.EmailItf
	redis          redis.RedisItf
	oauth          oauth.OAuthItf
}

func NewAuthUsecase(userRepository userRepository.UserRepositoryItf, db *gorm.DB, conf *conf.Config, jwt jwt.JWTItf, email email.EmailItf, redis redis.RedisItf, oauth oauth.OAuthItf) AuthUsecaseItf {
	return &AuthUsecase{
		userRepository: userRepository,
		db:             db,
		conf:           conf,
		jwt:            jwt,
		email:          email,
		redis:          redis,
		oauth:          oauth,
	}
}

func (uc *AuthUsecase) Register(req dto.RegisterRequest) *res.Err {
	user, err := uc.userRepository.GetUserByEmail(req.Email)
	if err != nil {
		return res.ErrInternalServerError(res.FailedFindUser)
	}

	needUpdatePassword := user != nil && user.GoogleID != nil && user.Password == nil

	if user != nil && !needUpdatePassword {
		return res.ErrConflict(res.EmailAlreadyExists)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return res.ErrInternalServerError(res.FailedHashPassword)
	}

	hashedPassword := string(hashed)

	if needUpdatePassword {
		if err := uc.userRepository.UpdateUser(req.Email, &entity.User{
			Password: &hashedPassword,
		}); err != nil {
			return res.ErrInternalServerError(res.FailedUpdateUser)
		}
	} else if user == nil {
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
	expiration := uc.conf.OTPExpiry

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

func (uc *AuthUsecase) GoogleLogin() (string, *res.Err) {
	stateLength := uc.conf.StateLength
	bytes := make([]byte, stateLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", res.ErrInternalServerError(res.FailedGenerateOAuthState)
	}

	state := base64.RawURLEncoding.EncodeToString(bytes)
	if len(state) > stateLength {
		state = state[:stateLength]
	}

	if err := uc.redis.SetOAuthState(state, []byte(state), uc.conf.StateExpiry); err != nil {
		return "", res.ErrInternalServerError(res.FailedStoreOAuthState)
	}

	url, err := uc.oauth.GenerateLink(state)
	if err != nil {
		return "", res.ErrInternalServerError(res.FailedGenerateOAuthLink)
	}

	return url, nil
}

func (uc *AuthUsecase) GoogleCallback(req *dto.GoogleCallbackRequest) (string, string, bool, *res.Err) {
	if req.Error != "" {
		return "", "", false, res.ErrInternalServerError(res.FailedOAuthCallback)
	}

	state, err := uc.redis.GetOAuthState(req.State)
	if err != nil {
		return "", "", false, res.ErrUnauthorized(res.OAuthStateNotFound)
	}

	if string(state) != req.State {
		return "", "", false, res.ErrUnauthorized(res.OAuthStateInvalid)
	}

	if err := uc.redis.DeleteOAuthState(req.State); err != nil {
		return "", "", false, res.ErrInternalServerError(res.FailedDeleteOAuthState)
	}

	token, err := uc.oauth.ExchangeToken(req.Code)
	if err != nil {
		return "", "", false, res.ErrInternalServerError(res.FailedExchangeOAuthToken)
	}

	profile, err := uc.oauth.GetProfile(token)
	if err != nil {
		return "", "", false, res.ErrInternalServerError(res.FailedGetOAuthProfile)
	}

	isNewUser := false

	user, err := uc.userRepository.GetUserByEmail(profile.Email)
	if err != nil {
		return "", "", false, res.ErrInternalServerError(res.FailedFindUser)
	}

	if user != nil {
		if user.GoogleID == nil {
			if err := uc.userRepository.UpdateUser(user.Email, &entity.User{
				GoogleID: &profile.ID,
			}); err != nil {
				return "", "", false, res.ErrInternalServerError(res.FailedUpdateUser)
			}
		}
	} else {
		isNewUser = true

		user = &entity.User{
			Name:     profile.Name,
			Email:    profile.Email,
			GoogleID: &profile.ID,
			Verified: profile.Verified,
		}

		if err := uc.userRepository.CreateUser(user); err != nil {
			return "", "", false, res.ErrInternalServerError(res.FailedCreateUser)
		}
	}

	if !user.Verified {
		return "", "", false, res.ErrUnauthorized(res.UserNotVerified)
	}

	refreshToken, err := uc.jwt.GenerateRefershToken(user.ID, false)
	if err != nil {
		return "", "", false, res.ErrInternalServerError(res.FailedGenerateRefreshToken)
	}

	refreshTokens, err := uc.userRepository.GetRefreshTokens(user.ID)
	if err != nil {
		return "", "", false, res.ErrInternalServerError(res.FailedGetRefreshTokens)
	}

	if len(refreshTokens) >= 2 {
		if err := uc.userRepository.RemoveRefreshToken(refreshTokens[0].Token); err != nil {
			return "", "", false, res.ErrInternalServerError(res.FailedRemoveRefreshToken)
		}
	}

	if err := uc.userRepository.AddRefreshToken(user.ID, refreshToken); err != nil {
		return "", "", false, res.ErrInternalServerError(res.FailedAddRefreshToken)
	}

	accessToken, err := uc.jwt.GenerateAccessToken(user.ID, user.Name, user.Email)
	if err != nil {
		return "", "", false, res.ErrInternalServerError(res.FailedGenerateAccessToken)
	}

	return accessToken, refreshToken, isNewUser, nil
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

func (uc *AuthUsecase) Logout(req dto.LogoutRequest) *res.Err {
	user, err := uc.userRepository.GetUserByRefreshToken(req.RefreshToken)
	if err != nil {
		return res.ErrInternalServerError(res.FailedFindUser)
	}

	if user == nil {
		return res.ErrUnauthorized(res.InvalidRefreshToken)
	}

	if err := uc.userRepository.RemoveRefreshToken(req.RefreshToken); err != nil {
		return res.ErrInternalServerError(res.FailedRemoveRefreshToken)
	}

	return nil
}
