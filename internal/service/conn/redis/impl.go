package redis

import (
	"fmt"
	"go-rate-limiter/deployment/config"
	"go-rate-limiter/internal/service/base"

	"github.com/go-redis/redis/v8"
)

type Impl struct {
	Client *redis.Client
}

func NewRedisClient(cfg *config.Config) RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
	})

	if _, err := client.Ping(base.Background()).Result(); err != nil {
		panic(err)
	}

	return &Impl{
		Client: client,
	}
}

func (i *Impl) Ping(ctx base.Ctx) error {
	if _, err := i.Client.Ping(ctx).Result(); err != nil {
		ctx.WithField("err", err).Error("Failed to RedisClient.Ping")
		return err
	}
	return nil
}

func (i *Impl) RunScript(ctx base.Ctx, script string, keys []string, args ...interface{}) (interface{}, error) {
	val, err := redis.NewScript(script).Run(ctx, i.Client, keys, args).Result()
	if err != nil && err != redis.Nil {
		ctx.WithField("err", err).Error("Failed to RedisClient.RunScript")
	}
	return val, nil
}

func (i *Impl) ClearAll(ctx base.Ctx) error {
	if _, err := i.Client.FlushAll(ctx).Result(); err != nil {
		ctx.WithField("err", err).Error("client.FlushAll failed")
		return err
	}
	return nil
}
