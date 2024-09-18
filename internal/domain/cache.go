package domain

import (
	"time"
)

type RedisCacheRepository interface {
	Get(key string) (string, error)
	Set(key string, value interface{}, expiration time.Duration) error
}
