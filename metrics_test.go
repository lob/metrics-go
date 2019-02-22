package metrics

import (
	"errors"
	"testing"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/lob/metrics-go/pkg/lambda"

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

type mockWriter struct {
}

func (w *mockWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (w *mockWriter) Close() error {
	return nil
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

func (m *mockClient) Gauge(name string, duration float64, tags []string, rate float64) error {
	m.name = testMetric
	m.duration = testDuration
	m.tags = testTags
	m.rate = testRate

	return errors.New("test error")
}

func (m *mockClient) Close() error {
	return errors.New("test error")
}

func newMockedClient(t *testing.T, cfg Config) *StatsReporter {
	return &StatsReporter{
		client: &mockClient{t, "", 0, 0, []string{}, 0},
	}
}

func TestNewMetrics(t *testing.T) {
	t.Run("create new statsd", func(t *testing.T) {
		cfg := Config{
			Namespace:  "testing.",
			StatsdHost: "127.0.0.1",
			StatsdPort: 8125,
		}
		m, err := New(cfg)
		assert.NoError(t, err)

		client, ok := m.client.(*statsd.Client)
		require.True(t, ok)
		assert.Equal(t, "testing.", client.Namespace)
	})

	t.Run("creates new lambda", func(t *testing.T) {
		cfg := Config{
			Namespace:    "testing.",
			Lambda:       true,
			LambdaLogger: &mockWriter{},
		}
		m, err := New(cfg)
		assert.NoError(t, err)

		client, ok := m.client.(*lambda.Client)
		require.True(t, ok)
		assert.Equal(t, "testing.", client.Namespace)
	})

	t.Run("namespace required", func(t *testing.T) {
		_, err := New(Config{
			StatsdHost: "127.0.0.1",
			StatsdPort: 8125,
		})
		assert.Error(t, err)
	})

	t.Run("namespace appends separator", func(t *testing.T) {
		m, err := New(Config{
			Namespace:  "testing",
			StatsdHost: "127.0.0.1",
			StatsdPort: 8125,
		})
		assert.NoError(t, err)

		client, ok := m.client.(*statsd.Client)
		require.True(t, ok)
		assert.Equal(t, "testing.", client.Namespace)
	})
}

func TestCount(t *testing.T) {
	cfg := Config{
		StatsdHost: "127.0.0.1",
		StatsdPort: 8125,
	}

	t.Run("calls Count function and ignores error", func(tt *testing.T) {
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

func TestGauge(t *testing.T) {
	cfg := Config{
		StatsdHost: "127.0.0.1",
		StatsdPort: 8125,
	}

	t.Run("calls Gauge function and ignores error", func(tt *testing.T) {
		metrics := newMockedClient(t, cfg)

		metrics.Gauge(testMetric, testDuration, testTags...)

		mc, ok := metrics.client.(*mockClient)
		require.True(t, ok, "unexpected error during type assertion")

		assert.Equal(t, testMetric, mc.name, "inconsistent metric name")
		assert.Equal(t, testDuration, mc.duration, "inconsistent metric duration")
		assert.Equal(t, testTags, mc.tags, "inconsistent tags")
		assert.Equal(t, testRate, mc.rate, "inconsistent rate")
	})
}

func TestHistogram(t *testing.T) {
	cfg := Config{
		StatsdHost: "127.0.0.1",
		StatsdPort: 8125,
	}

	t.Run("calls Histogram function and ignores error", func(tt *testing.T) {
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
