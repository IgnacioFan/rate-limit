package ratelimiter

import "go-rate-limiter/internal/service/base"

type Ratelimiter interface {
	AcquireByIP(ctx base.Ctx, key string) (permit bool, count uint)
}
