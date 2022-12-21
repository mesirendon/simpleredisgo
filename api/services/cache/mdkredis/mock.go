package mdkredis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type ServiceMock struct {
	redis.Cmdable
	records map[string]map[string]any
}

// Redis mocked methods

func (m *ServiceMock) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return m.records["Set"][key].(*redis.StatusCmd)
}

func (m *ServiceMock) Get(ctx context.Context, key string) *redis.StringCmd {
	return m.records["Get"][key].(*redis.StringCmd)
}

func (m *ServiceMock) TTL(ctx context.Context, key string) *redis.DurationCmd {
	return m.records["TTL"][key].(*redis.DurationCmd)
}

func (m *ServiceMock) Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {
	return m.records["Scan"][match].(*redis.ScanCmd)
}

func (m *ServiceMock) FlushDB(ctx context.Context) *redis.StatusCmd {
	return m.records["FlushDB"]["FlushDB"].(*redis.StatusCmd)
}

func (m *ServiceMock) On(fnName string, key string, result any) {
	if m.records == nil {
		m.records = make(map[string]map[string]any)
	}

	if m.records[fnName] == nil {
		m.records[fnName] = make(map[string]any)
	}

	m.records[fnName][key] = result
}
