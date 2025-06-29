package redis

import (
	"encoding/json"
	"time"

	conf "github.com/Ablebil/sea-catering-be/config"
	"github.com/gofiber/storage/redis"
)

type RedisItf interface {
	SetCache(key string, data interface{}, exp time.Duration) error
	GetCache(key string, data interface{}) error
	DeleteCache(key string) error
	SetOTP(email string, otp string, exp time.Duration) error
	GetOTP(email string) (string, error)
	DeleteOTP(email string) error
	SetOAuthState(state string, value []byte, exp time.Duration) error
	GetOAuthState(state string) ([]byte, error)
	DeleteOAuthState(state string) error
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

func (r *Redis) SetCache(key string, data interface{}, exp time.Duration) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return r.store.Set(key, bytes, exp)
}

func (r *Redis) GetCache(key string, data interface{}) error {
	val, err := r.store.Get(key)
	if err != nil {
		return err
	}

	return json.Unmarshal(val, data)
}

func (r *Redis) DeleteCache(key string) error {
	return r.store.Delete(key)
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

func (r *Redis) SetOAuthState(state string, value []byte, exp time.Duration) error {
	key := "gstate:" + state
	return r.store.Set(key, value, exp)
}

func (r *Redis) GetOAuthState(state string) ([]byte, error) {
	key := "gstate:" + state
	val, err := r.store.Get(key)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (r *Redis) DeleteOAuthState(state string) error {
	key := "gstate:" + state
	return r.store.Delete(key)
}
