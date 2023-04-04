package redis

import (
	"go-rate-limiter/internal/service/base"

	"github.com/go-redis/redis/v8"
)

type Impl struct {
	Client *redis.Client
}

func NewRedisClient() RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
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
		ctx.WithField("err", err).Error("client.Ping failed")
		return err
	}
	return nil
}
