package ratelimiter

type Ratelimiter interface {
	AcquireByIP(key string) (permit bool, count uint)
}
