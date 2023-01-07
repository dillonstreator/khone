package khone

import (
	"context"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type FlowRate struct {
	Timestamp time.Time `json:"timestamp"`
	GPM       float32   `json:"gpm"`
}

func (f *FlowRate) UnmarshalCSV(m map[string]string) error {
	timestamp, err := time.Parse(time.RFC3339, m["timestamp"])
	if err != nil {
		return err
	}
	f.Timestamp = timestamp

	gpm, err := strconv.ParseFloat(m["gpm"], 32)
	if err != nil {
		return err
	}
	f.GPM = float32(gpm)

	return nil
}

func TestStreamCSV(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	r := strings.NewReader(`timestamp,gpm
2023-01-13T03:03:00Z,1.4
2023-01-13T03:04:00Z,1.9
2023-01-13T03:05:00Z,1.1`)

	sleepDuration := time.Millisecond * 100
	startTime := time.Now()

	var calls int32

	err := StreamCSV(ctx, r, func(fr *FlowRate, i int) {
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

func TestStreamCSV_context(t *testing.T) {
	assert := assert.New(t)

	sleepDuration := time.Millisecond * 100

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*90)
	defer cancel()

	r := strings.NewReader(`timestamp,gpm
2023-01-13T03:03:00Z,1.4
2023-01-13T03:04:00Z,1.9
2023-01-13T03:05:00Z,1.1`)

	var calls int32

	err := StreamCSV(ctx, r, func(fr *FlowRate, i int) {
		atomic.AddInt32(&calls, 1)

		assert.Equal("2023-01-13T03:03:00Z", fr.Timestamp.Format(time.RFC3339))
		assert.Equal(float32(1.4), fr.GPM)

		time.Sleep(sleepDuration)
	})
	assert.ErrorIs(err, context.DeadlineExceeded)

	assert.Equal(int32(1), calls)
}
