package strategy

import "go-rate-limiter/internal/service/base"

type Strategy interface {
	Acquire(context base.Ctx, key string) (bool, int, error)
}
