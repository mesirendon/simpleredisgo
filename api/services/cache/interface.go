package cache

import (
	"time"
)

type ICache interface {
	Set(key string, value any, expiration time.Duration) error
	Get(key string, output any) error
	GetRecordsByPattern(regexp string, output any) error
	TTL(key string) (time.Duration, error)
	KeyPatternRecordSize(regexp string) (int, error)
	Close() error
	FlushDB() error
}
