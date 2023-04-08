package tokenbucket

import (
	"fmt"
	"go-rate-limiter/internal/service/base"
	"go-rate-limiter/internal/service/conn/redis"
	"go-rate-limiter/internal/service/ratelimiter/strategy"
	"time"
)

const (
	script string = `
	local capacity = tonumber(ARGV[4])
	local refill = tonumber(ARGV[3])
	local refillToken = tonumber(ARGV[2])
	local ts = tonumber(ARGV[1])
	local lastUpdate = ts
	local remainToken = capacity

	local last = redis.call('HMGET', KEYS[1], 'ts', 'tokens')
	if last[1] then
		local lastTs = tonumber(last[1])
		local lastTokens = tonumber(last[2])
		local refillCount = math.floor((ts - lastTs) / refill)


		remainToken = math.min(capacity, lastTokens + (refillCount * refillToken))
		lastUpdate = math.min(ts, lastTs + (refillCount * refill))
	end

	if remainToken >= 0 then
			remainToken = remainToken - 1
	end
	redis.call('HMSET', KEYS[1], 'ts', ts, 'tokens', remainToken)
	redis.call('EXPIRE', KEYS[1], math.ceil(capacity / refill))
	return remainToken
	`
)

var (
	timeNow = time.Now

	bucketSize      int     = 10
	refillPerSecond float64 = 60
	refillToken     int     = 10
)

type Impl struct {
	Client      redis.RedisClient
	Size        int
	Refill      float64
	RefillToken int
}

func NewTokenBucket(client redis.RedisClient) strategy.Strategy {
	return &Impl{
		Client:      client,
		Size:        bucketSize,
		Refill:      refillPerSecond,
		RefillToken: refillToken,
	}
}

func (i *Impl) Acquire(context base.Ctx, key string) (bool, int, error) {
	now := timeNow()
	cacheKey := fmt.Sprintf("tokenbucket:%s", key)

	remain, err := i.Client.RunScript(
		context,
		script,
		[]string{cacheKey},
		now.Unix(),
		i.RefillToken,
		i.Refill,
		i.Size,
	)
	if err != nil {
		context.WithField("err", err).Error("redis.RunScript failed")
		return false, 0, err
	}
	if remain.(int64) < 0 {
		return false, i.Size, nil
	}
	return true, int(remain.(int64)), nil
}
