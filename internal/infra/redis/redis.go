package redis

import (
	"time"

	conf "github.com/Ablebil/sea-catering-be/config"
	"github.com/gofiber/storage/redis"
)

type RedisItf interface {
	SetOTP(email string, otp string, exp time.Duration) error
	GetOTP(email string) (string, error)
	DeleteOTP(email string) error
}

type Redis struct {
	store *redis.Storage
}

func NewRedis(conf *conf.Config) RedisItf {
	return &Redis{
		store: redis.New(redis.Config{
			Host:     conf.RedisHost,
			Port:     conf.RedisPort,
			Password: conf.RedisPassword,
		}),
	}
}

func (r *Redis) SetOTP(email string, otp string, exp time.Duration) error {
	key := "otp:" + email
	return r.store.Set(key, []byte(otp), exp)
}

func (r *Redis) GetOTP(email string) (string, error) {
	key := "otp:" + email
	val, err := r.store.Get(key)
	if err != nil {
		return "", err
	}

	return string(val), nil
}

func (r *Redis) DeleteOTP(email string) error {
	key := "otp:" + email
	return r.store.Delete(key)
}
