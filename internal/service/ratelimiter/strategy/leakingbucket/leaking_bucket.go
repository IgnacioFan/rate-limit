package leakingbucket

import (
	"go-rate-limiter/internal/service/base"
	"go-rate-limiter/internal/service/conn/redis"
	"go-rate-limiter/internal/service/ratelimiter/strategy"
	"time"
)

// Copyright (c) 2017 Pavel Pravosud
// https://github.com/rwz/redis-gcra/blob/master/vendor/perform_gcra_ratelimit.lua
const (
	script = `
		local burst = ARGV[1]
		local rate = ARGV[2]
		local period = ARGV[3]
		local cost = ARGV[4]

		local emissionInterval = period / rate
		local increment = emissionInterval * cost
		local burstOffset = emissionInterval * burst

		local now = tonumber(ARGV[5])
		local tat = redis.call("GET", KEYS[1])
		if not tat then
			tat = now
		else
			tat = tonumber(tat)
		end

		local newTat = math.max(tat, now) + increment
		local allowAt = newTat - burstOffset

		local diff = now - allowAt
		local remaining = math.floor(diff / emissionInterval)

		if remaining < 0 then
			return remaining
		end

		local resetAfter = newTat - now
		if resetAfter > 0 then
			redis.call("SET", KEYS[1], newTat, "EX", math.ceil(resetAfter))
		end

		return remaining
	`
)

var (
	timeNow = time.Now

	burst  int           = 1000
	rate   int           = 100
	period time.Duration = 60
	cost   int           = 20
)

type Impl struct {
	Client redis.RedisClient
	Burst  int
	Rate   int
	Period time.Duration
	Cost   int
}

func NewLeakingBucket(client redis.RedisClient) strategy.Strategy {
	return &Impl{
		Client: client,
		Burst:  burst,
		Rate:   rate,
		Period: period,
		Cost:   cost,
	}
}

func (i *Impl) Acquire(ctx base.Ctx, key string) (bool, int, error) {
	now := timeNow()

	remain, err := i.Client.RunScript(
		ctx,
		script,
		[]string{key},
		i.Burst,
		i.Rate,
		i.Period,
		i.Cost,
		now.Unix(),
	)
	if err != nil {
		ctx.WithField("err", err).Error("Redis.Runscript failed")
		return false, 0, err
	}
	if remain.(int64) < 0 {
		return false, i.Burst, nil
	}
	return true, int(remain.(int64)), nil
}
