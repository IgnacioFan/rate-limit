package leakingbucket

import (
	"go-rate-limiter/internal/service/base"
	"go-rate-limiter/internal/service/conn/redis"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	client  = redis.NewRedisClient()
	mockCtx = base.Background()
	ipAddr  = "localhost"
	mockNow = time.Date(2023, time.April, 1, 0, 0, 0, 0, time.UTC)
)

type mockFuncs struct {
	mock.Mock
}

func (m *mockFuncs) timeNow() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

func TestLeakingBucketAcquire(t *testing.T) {
	leakingbucket := NewLeakingBucket(client).(*Impl)
	leakingbucket.Burst = 50
	leakingbucket.Rate = 5
	leakingbucket.Period = 3
	leakingbucket.Cost = 10

	// mock timeNow function
	mockFuncs := new(mockFuncs)
	timeNow = mockFuncs.timeNow

	tests := []struct {
		Desc        string
		AcquireTime []time.Time
		Exp         []bool
		ExpCount    []int
	}{
		{
			"Get throttled",
			[]time.Time{
				mockNow,
				mockNow.Add(1 * time.Second),
				mockNow.Add(2 * time.Second),
				mockNow.Add(3 * time.Second),
				mockNow.Add(4 * time.Second),
				mockNow.Add(5 * time.Second),
			},
			[]bool{true, true, true, true, true, false},
			[]int{40, 31, 23, 15, 6, 50},
		},
		{
			"reset rate limit per 3 secs",
			[]time.Time{
				mockNow,
				mockNow.Add(1 * time.Second),
				mockNow.Add(12 * time.Second),
			},
			[]bool{true, true, true},
			[]int{40, 31, 40},
		},
	}
	for _, test := range tests {
		for i, slot := range test.AcquireTime {
			mockFuncs.On("timeNow").Return(slot).Once()

			permit, remain, err := leakingbucket.Acquire(mockCtx, ipAddr)
			t.Run(test.Desc, func(t *testing.T) {
				assert.Equal(t, err, nil)
				assert.Equal(t, test.Exp[i], permit)
				assert.Equal(t, test.ExpCount[i], remain)
			})
		}
		client.ClearAll(mockCtx)
	}
}
