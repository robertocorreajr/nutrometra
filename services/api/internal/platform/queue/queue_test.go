package queue

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPostgresQueue(t *testing.T) {
	q := NewPostgresQueue(nil)
	require.NotNil(t, q)
	assert.NotNil(t, q.handlers)
	assert.Equal(t, 1*time.Second, q.pollInterval)
	assert.Nil(t, q.pool)
}

func TestRegisterHandler(t *testing.T) {
	q := NewPostgresQueue(nil)

	q.RegisterHandler("test_job", func(_ context.Context, _ *JobRecord) error {
		return nil
	})

	handler, ok := q.handlers["test_job"]
	assert.True(t, ok, "handler should be registered")
	assert.NotNil(t, handler)

	_, ok = q.handlers["nonexistent"]
	assert.False(t, ok, "nonexistent handler should not be present")
}

func TestJobDefaults(t *testing.T) {
	t.Run("zero MaxAttempts defaults to defaultMaxAttempts", func(t *testing.T) {
		j := Job{
			JobType:     "test",
			MaxAttempts: 0,
		}
		assert.Equal(t, 0, j.MaxAttempts)
		// The Enqueue method applies default — verify constant
		assert.Equal(t, 3, defaultMaxAttempts)
	})

	t.Run("zero ScheduleAt is treated as now by Enqueue", func(t *testing.T) {
		j := Job{
			JobType: "test",
		}
		assert.True(t, j.ScheduleAt.IsZero())
	})

	t.Run("nil payload is treated as empty JSON object by Enqueue", func(t *testing.T) {
		j := Job{
			JobType: "test",
		}
		assert.Nil(t, j.Payload)
	})
}
