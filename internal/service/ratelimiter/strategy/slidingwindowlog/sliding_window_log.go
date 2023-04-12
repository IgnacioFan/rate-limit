package slidingwindowlog

import (
	"fmt"
	"go-rate-limiter/internal/service/base"
	"go-rate-limiter/internal/service/conn/redis"
	"go-rate-limiter/internal/service/ratelimiter/strategy"
	"time"
)

const (
	script string = `
	local size = tonumber(ARGV[1])
	local limitPerSec = tonumber(ARGV[2])
	local now = tonumber(ARGV[3])

	local min = now - limitPerSec
	local max = now

	local count = redis.call('ZCOUNT', KEYS[1], min, max)

	redis.call('ZADD', KEYS[1], now, max)

	if count == 0 then
		-- clear the key
		redis.call('EXPIRE', KEYS[1], limitPerSec)
		-- clear previous window
		redis.call('ZREMRANGEBYSCORE', KEYS[1], '-inf', min)
	end
	return size - count - 1
	`
)

var (
	timeNow = time.Now

	Size        = 60
	LimitPerSec = 60
)

type Impl struct {
	Client      redis.RedisClient
	Size        int
	LimitPerSec int
}

func NewSlidingWindowLog(client redis.RedisClient) strategy.Strategy {
	return &Impl{
		Client:      client,
		Size:        Size,
		LimitPerSec: LimitPerSec,
	}
}

func (i *Impl) Acquire(ctx base.Ctx, key string) (bool, int, error) {
	now := timeNow()
	redisKey := fmt.Sprintf("sliding_window_log:%s", key)
	remain, err := i.Client.RunScript(ctx, script, []string{redisKey}, i.Size, i.LimitPerSec, now.Unix())
	if err != nil {
		ctx.WithField("err", err).Error("redis.RunScript failed")
		return false, 0, err
	}

	if int(remain.(int64)) < 0 {
		return false, i.Size, nil
	}
	return true, int(remain.(int64)), nil
}
