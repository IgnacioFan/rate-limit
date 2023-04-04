package redis

import "go-rate-limiter/internal/service/base"

type RedisClient interface {
	Ping(ctx base.Ctx) error
}
