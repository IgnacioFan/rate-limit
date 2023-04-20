package fixedwindow

import (
	"go-rate-limiter/deployment/config"
	"go-rate-limiter/internal/service/base"
	"go-rate-limiter/internal/service/conn/redis"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	client  = redis.NewRedisClient(config.NewConfig())
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

func TestFixedWindowAcquire(t *testing.T) {
	fixedWindow := NewFixedWindow(client).(*Impl)
	fixedWindow.Size = 5
	fixedWindow.LimitPerSec = 10

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
				mockNow.Add(6 * time.Second),
			},
			[]bool{true, true, true, true, true, false, false},
			[]int{4, 3, 2, 1, 0, 5, 5},
		},
		{
			"reset request in 10 secs",
			[]time.Time{
				mockNow,
				mockNow.Add(1 * time.Second),
				mockNow.Add(10 * time.Second),
			},
			[]bool{true, true, true},
			[]int{4, 3, 4},
		},
	}
	for _, test := range tests {
		for i, slot := range test.AcquireTime {
			mockFuncs.On("timeNow").Return(slot).Once()

			permit, remain, err := fixedWindow.Acquire(mockCtx, ipAddr)
			t.Run(test.Desc, func(t *testing.T) {
				assert.Equal(t, err, nil)
				assert.Equal(t, test.Exp[i], permit)
				assert.Equal(t, test.ExpCount[i], remain)
			})
		}
		client.ClearAll(mockCtx)
	}
}
