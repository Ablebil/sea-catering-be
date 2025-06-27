package jwt

import (
	"errors"
	"time"

	conf "github.com/Ablebil/sea-catering-be/config"
	"github.com/Ablebil/sea-catering-be/internal/domain/entity"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTItf interface {
	GenerateAccessToken(userId uuid.UUID, name string, email string, role entity.UserRole) (string, error)
	GenerateRefershToken(userId uuid.UUID, rememberMe bool) (string, error)
	VerifyAccessToken(token string) (uuid.UUID, string, string, error)
	VerifyRefreshToken(token string) (uuid.UUID, error)
}

type JWT struct {
	accessSecret  string
	refreshSecret string
}

func NewJWT(conf *conf.Config) JWTItf {
	return &JWT{
		accessSecret:  conf.AccessSecret,
		refreshSecret: conf.RefreshSecret,
	}
}

type AccessClaims struct {
	UserID uuid.UUID       `json:"user_id"`
	Name   string          `json:"name"`
	Email  string          `json:"email"`
	Role   entity.UserRole `json:"role"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID     uuid.UUID `json:"user_id"`
	RememberMe bool      `json:"remember_me"`
	jwt.RegisteredClaims
}

func (j *JWT) GenerateAccessToken(userId uuid.UUID, name string, email string, role entity.UserRole) (string, error) {
	claims := AccessClaims{
		UserID: userId,
		Name:   name,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.accessSecret))
}

func (j *JWT) GenerateRefershToken(userId uuid.UUID, rememberMe bool) (string, error) {
	var ttl time.Duration
	if rememberMe {
		ttl = 30 * 24 * time.Hour
	} else {
		ttl = 7 * 24 * time.Hour
	}

	claims := RefreshClaims{
		UserID:     userId,
		RememberMe: rememberMe,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.refreshSecret))
}

func (j *JWT) VerifyAccessToken(tokenString string) (uuid.UUID, string, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.accessSecret), nil
	})

	if err != nil {
		return uuid.Nil, "", "", err
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid {
		return uuid.Nil, "", "", errors.New("couldn't parse access token claims")
	}

	return claims.UserID, claims.Name, claims.Email, nil
}

func (j *JWT) VerifyRefreshToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.refreshSecret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid {
		return uuid.Nil, errors.New("couldn't parse refresh token claims")
	}

	return claims.UserID, nil
}
