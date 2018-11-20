package metrics

import (
	"errors"
	"testing"

	"github.com/lob/assets-proxy/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testCount = int64(1)
const testDuration = float64(50)
const testRate = float64(1)
const testMetric = "test_metric"
const testTag = "foo:bar"

var testTags = []string{testTag}

type mockClient struct {
	t        *testing.T
	name     string
	count    int64
	duration float64
	tags     []string
	rate     float64
}

func (m *mockClient) Count(name string, count int64, tags []string, rate float64) error {
	m.name = testMetric
	m.count = testCount
	m.tags = testTags
	m.rate = testRate

	return errors.New("test error")
}

func (m *mockClient) Histogram(name string, duration float64, tags []string, rate float64) error {
	m.name = testMetric
	m.duration = testDuration
	m.tags = testTags
	m.rate = testRate

	return errors.New("test error")
}

func newMockedClient(t *testing.T, cfg config.Config) Metrics {
	metrics, err := New(cfg)
	require.NoError(t, err)
	require.NotNil(t, metrics)

	metrics.client = &mockClient{t, "", 0, 0, []string{}, 0}

	return metrics
}

func TestCount(t *testing.T) {
	cfg, err := config.New()
	require.NoError(t, err)

	t.Run("calls Datadog Count function and ignores error", func(tt *testing.T) {
		metrics := newMockedClient(t, cfg)

		metrics.Count(testMetric, testCount, testTags...)

		mc, ok := metrics.client.(*mockClient)
		require.True(t, ok, "unexpected error during type assertion")

		assert.Equal(t, testMetric, mc.name, "inconsistent metric name")
		assert.Equal(t, testCount, mc.count, "inconsistent metric count")
		assert.Equal(t, testTags, mc.tags, "inconsistent tags")
		assert.Equal(t, testRate, mc.rate, "inconsistent rate")
	})
}

func TestHistogram(t *testing.T) {
	cfg, err := config.New()
	require.NoError(t, err)

	t.Run("calls Datadog Histogram function and ignores error", func(tt *testing.T) {
		metrics := newMockedClient(t, cfg)

		metrics.Histogram(testMetric, testDuration, testTags...)

		mc, ok := metrics.client.(*mockClient)
		require.True(t, ok, "unexpected error during type assertion")

		assert.Equal(t, testMetric, mc.name, "inconsistent metric name")
		assert.Equal(t, testDuration, mc.duration, "inconsistent duration")
		assert.Equal(t, testTags, mc.tags, "inconsistent tags")
		assert.Equal(t, testRate, mc.rate, "inconsistent rate")
	})
}
