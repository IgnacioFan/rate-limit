package fixedwindow

import (
	"fmt"
	"go-rate-limiter/internal/service/base"
	"go-rate-limiter/internal/service/conn/redis"
	"go-rate-limiter/internal/service/ratelimiter/strategy"
	"time"
)

const (
	script string = `
		local val = redis.call('INCR', KEYS[1])
		if val == 1 then
			redis.call('EXPIRE', KEYS[1], ARGV[2])
		end
		return ARGV[1] - val
	`
)

var (
	timeNow = time.Now

	fixedWindowSize        int = 60
	fixedWindowLimitPerSec int = 60
)

type Impl struct {
	Client      redis.RedisClient
	Size        int
	LimitPerSec int
}

func NewFixedWindow(client redis.RedisClient) strategy.Strategy {
	return &Impl{
		Client:      client,
		Size:        fixedWindowSize,
		LimitPerSec: fixedWindowLimitPerSec,
	}
}

func (i *Impl) Acquire(ctx base.Ctx, key string) (bool, int, error) {
	now := timeNow()
	window := now.Unix() / int64(i.LimitPerSec)
	redisKey := fmt.Sprintf("fixed_window:%v:%v", key, window)
	remain, err := i.Client.RunScript(ctx, script, []string{redisKey}, i.Size, i.LimitPerSec)
	if err != nil {
		ctx.WithField("err", err).Error("redis.RunScript failed")
		return false, 0, err
	}
	if remain.(int64) < 0 {
		return false, i.Size, nil
	}
	return true, int(remain.(int64)), nil
}
