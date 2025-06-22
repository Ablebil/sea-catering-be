package jwt

import (
	"errors"
	"time"

	conf "github.com/Ablebil/sea-catering-be/config"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTItf interface {
	GenerateAccessToken(userId uuid.UUID, name string, email string) (string, error)
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
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID     uuid.UUID `json:"user_id"`
	RememberMe bool      `json:"remember_me"`
	jwt.RegisteredClaims
}

func (j *JWT) GenerateAccessToken(userId uuid.UUID, name string, email string) (string, error) {
	claims := AccessClaims{
		UserID: userId,
		Name:   name,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * 60)),
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
	return token.SignedString([]byte(j.accessSecret))
}

func (j *JWT) VerifyAccessToken(tokenString string) (uuid.UUID, string, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.accessSecret), nil
	})

	if err != nil || !token.Valid {
		return uuid.Nil, "", "", errors.New("Invalid access token")
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok {
		return uuid.Nil, "", "", errors.New("Couldn't parse access token claims")
	}

	return claims.UserID, claims.Name, claims.Email, nil
}

func (j *JWT) VerifyRefreshToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.refreshSecret), nil
	})

	if err != nil || !token.Valid {
		return uuid.Nil, errors.New("Invalid refresh token")
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok {
		return uuid.Nil, errors.New("Couldn't parse refresh token claims")
	}

	return claims.UserID, nil
}
