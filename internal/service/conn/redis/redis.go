package redis

import "go-rate-limiter/internal/service/base"

type RedisClient interface {
	Ping(ctx base.Ctx) error
	RunScript(ctx base.Ctx, script string, keys []string, args ...interface{}) (interface{}, error)
	ClearAll(ctx base.Ctx) error
}
