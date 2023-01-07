package khone

import (
	"context"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStreamJSON(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	r := strings.NewReader(`[{"timestamp":"2023-01-13T03:03:00Z","gpm":1.4},{"timestamp":"2023-01-13T03:04:00Z","gpm":1.9},{"timestamp":"2023-01-13T03:05:00Z","gpm":1.1}]`)

	sleepDuration := time.Millisecond * 100
	startTime := time.Now()

	var calls int32

	err := StreamJSON(ctx, r, func(fr *FlowRate, i int) {
		atomic.AddInt32(&calls, 1)

		assert.WithinDuration(startTime, time.Now(), sleepDuration)

		switch i {
		case 0:
			assert.Equal("2023-01-13T03:03:00Z", fr.Timestamp.Format(time.RFC3339))
			assert.Equal(float32(1.4), fr.GPM)
		case 1:
			assert.Equal("2023-01-13T03:04:00Z", fr.Timestamp.Format(time.RFC3339))
			assert.Equal(float32(1.9), fr.GPM)
		case 2:
			assert.Equal("2023-01-13T03:05:00Z", fr.Timestamp.Format(time.RFC3339))
			assert.Equal(float32(1.1), fr.GPM)
		}

		time.Sleep(sleepDuration)
	}, WithConcurrency(3))
	assert.NoError(err)

	assert.Equal(int32(3), calls)
}

func TestStreamJSON_context(t *testing.T) {
	assert := assert.New(t)

	sleepDuration := time.Millisecond * 100

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*90)
	defer cancel()

	r := strings.NewReader(`[{"timestamp":"2023-01-13T03:03:00Z","gpm":1.4},{"timestamp":"2023-01-13T03:04:00Z","gpm":1.9},{"timestamp":"2023-01-13T03:05:00Z","gpm":1.1}]`)

	var calls int32

	err := StreamJSON(ctx, r, func(fr *FlowRate, i int) {
		atomic.AddInt32(&calls, 1)

		assert.Equal("2023-01-13T03:03:00Z", fr.Timestamp.Format(time.RFC3339))
		assert.Equal(float32(1.4), fr.GPM)

		time.Sleep(sleepDuration)
	})
	assert.ErrorIs(err, context.DeadlineExceeded)

	assert.Equal(int32(1), calls)
}
