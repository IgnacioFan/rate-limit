package ratelimiter

import (
	"flag"
	"fmt"
	"go-rate-limiter/internal/service/base"
	"go-rate-limiter/internal/service/conn/redis"
	"go-rate-limiter/internal/service/ratelimiter/strategy"
	"go-rate-limiter/internal/service/ratelimiter/strategy/tokenbucket"
)

var (
	rate_limit_strategy = flag.String("rate_limit_strategy", "tokenbucket", "set strategy for rate limit")
)

type Impl struct {
	Strategy strategy.Strategy
}

func NewRatelimiter(client redis.RedisClient) Ratelimiter {
	var strategy strategy.Strategy
	switch *rate_limit_strategy {
	case "tokenbucket":
		strategy = tokenbucket.NewTokenBucket(client)
	case "leakingbucket":
		strategy = nil
		fmt.Println("NewLeakingBucket")
	case "slidingwindow":
		strategy = nil
		fmt.Println("NewSlidingWindow")
	case "fixedwindow":
		strategy = nil
		fmt.Println("NewFixedWindow")
	default:
		strategy = tokenbucket.NewTokenBucket(client)
	}
	return &Impl{
		Strategy: strategy,
	}
}

func (i *Impl) AcquireByIP(ctx base.Ctx, key string) (permit bool, count uint, err error) {
	fmt.Println("AcquireByIP", i.Strategy, *rate_limit_strategy, key, ctx)
	permit, remain, err := i.Strategy.Acquire(ctx, key)
	if err != nil {
		ctx.WithField("err", err).Error("strategy.acquire failed")
		return false, 0, err
	}
	return permit, uint(remain), nil
}
