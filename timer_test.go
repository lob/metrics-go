package metrics

import (
	"testing"
	"time"

	"github.com/lob/assets-proxy/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimer(t *testing.T) {
	cfg, err := config.New()
	require.NoError(t, err)

	t.Run("calls histogram with correct duration", func(tt *testing.T) {
		metrics := newMockedClient(t, cfg)

		timer := metrics.NewTimer(testMetric, testTag)
		require.NotNil(t, timer)

		time.Sleep(time.Duration(testDuration) * time.Millisecond)
		timer.End()

		mc, ok := metrics.client.(*mockClient)
		require.True(t, ok, "unexpected error during type assertion")

		assert.Equal(t, testMetric, mc.name, "inconsistent metric name")
		assert.True(t, testDuration <= mc.duration, "incorrect duration")
		assert.Equal(t, testTags, mc.tags, "inconsistent tags")
		assert.Equal(t, testRate, mc.rate, "inconsistent rate")
	})
}
