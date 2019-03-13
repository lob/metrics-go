package lambda

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testMetric = "test_metric"
const testCount = int64(1)
const testValue = float64(50)
const testRate = float64(1)
const testType = "test_type"
const testTag = "foo:bar"

var testTags = []string{testTag}

type testWriteCloser struct {
	t      *testing.T
	buffer *bytes.Buffer
	closed bool
}

func (w *testWriteCloser) Write(p []byte) (n int, err error) {
	return w.buffer.Write(p)
}

func (w *testWriteCloser) Close() error {
	w.closed = true
	return nil
}

func newTestClient(t *testing.T) *Client {
	w := &testWriteCloser{t, new(bytes.Buffer), false}
	c, err := New(w)
	assert.NoError(t, err)

	c.Namespace = "testing."
	return c
}

func TestNewLambda(t *testing.T) {
	w := &testWriteCloser{t, new(bytes.Buffer), false}

	t.Run("create new lambda Client", func(t *testing.T) {
		c, err := New(w)
		assert.NoError(t, err)

		assert.Equal(t, w, c.writer)
	})

	t.Run("errors for nil WriteCloser", func(t *testing.T) {
		_, err := New(nil)

		assert.Error(t, err)
	})
}

func TestClose(t *testing.T) {
	t.Run("calls Close function and closes the WriteCloser", func(t *testing.T) {
		tc := newTestClient(t)

		err := tc.Close()
		require.NoError(t, err)

		w := tc.writer.(*testWriteCloser)
		assert.Equal(t, true, w.closed)
	})
}

func TestCount(t *testing.T) {
	t.Run("calls Count function and calls send", func(t *testing.T) {
		tc := newTestClient(t)
		tc.Namespace = "test"

		err := tc.Count(testMetric, testCount, testTags, testRate)
		require.NoError(t, err)

		w := tc.writer.(*testWriteCloser)

		got := w.buffer.String()
		want := fmt.Sprintf(
			"MONITORING|%d|%s|%s|%s%s|#%s",
			time.Now().Unix(),
			strconv.FormatInt(testCount, 10),
			"count",
			tc.Namespace,
			testMetric,
			strings.Join(testTags, ","),
		)
		assert.Equal(t, got, want)
	})
}

func TestGauge(t *testing.T) {
	t.Run("calls Gauge function and calls send", func(t *testing.T) {
		tc := newTestClient(t)
		tc.Namespace = "test"

		err := tc.Gauge(testMetric, testValue, testTags, testRate)
		require.NoError(t, err)

		w := tc.writer.(*testWriteCloser)

		got := w.buffer.String()
		want := fmt.Sprintf(
			"MONITORING|%d|%s|%s|%s%s|#%s",
			time.Now().Unix(),
			strconv.FormatFloat(testValue, 'f', -1, 64),
			"gauge",
			tc.Namespace,
			testMetric,
			strings.Join(testTags, ","),
		)
		assert.Equal(t, got, want)
	})
}

func TestHistogram(t *testing.T) {
	t.Run("calls Histogram function and calls send", func(t *testing.T) {
		tc := newTestClient(t)
		tc.Namespace = "test"

		err := tc.Histogram(testMetric, testValue, testTags, testRate)
		require.NoError(t, err)

		w := tc.writer.(*testWriteCloser)

		got := w.buffer.String()
		want := fmt.Sprintf(
			"MONITORING|%d|%s|%s|%s%s|#%s",
			time.Now().Unix(),
			strconv.FormatFloat(testValue, 'f', -1, 64),
			"histogram",
			tc.Namespace,
			testMetric,
			strings.Join(testTags, ","),
		)
		assert.Equal(t, got, want)
	})
}

func TestSend(t *testing.T) {
	w := &testWriteCloser{t, new(bytes.Buffer), false}

	t.Run("calls send function and writes a string in the correct format", func(t *testing.T) {
		c, err := New(w)
		require.NoError(t, err)

		c.Namespace = "test"
		now := time.Now().Unix()
		name := "name"
		value := "value"
		tags := []string{"test:test", "other:other"}
		metricType := "type"

		err = c.send(name, value, tags, metricType)
		assert.NoError(t, err)

		got := w.buffer.String()
		want := fmt.Sprintf(
			"MONITORING|%d|%s|%s|%s%s|#%s",
			now,
			value,
			metricType,
			c.Namespace,
			name,
			strings.Join(tags, ","),
		)
		assert.Equal(t, got, want)
	})
}
