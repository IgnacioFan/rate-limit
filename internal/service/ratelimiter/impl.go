package ratelimiter

import (
	"flag"
	"fmt"
)

var (
	rate_limit_strategy = flag.String("rate_limit_strategy", "tokenbucket", "set strategy for rate limit")
)

type Impl struct {
	Strategy string
}

func NewRatelimiter() *Impl {
	var res string
	switch *rate_limit_strategy {
	case "tokenbucket":
		res = "NewTokenBucket"
	case "leakingbucket":
		res = "NewLeakingBucket"
	case "slidingwindow":
		res = "NewSlidingWindow"
	case "fixedwindow":
		res = "NewTokenBucket"
	default:
		res = "NewTokenBucket"
	}
	return &Impl{
		Strategy: res,
	}
}

func (i *Impl) AcquireByIP(key string) (permit bool, count uint) {
	fmt.Println("AcquireByIP", i.Strategy, *rate_limit_strategy, key)
	return false, 2
}
