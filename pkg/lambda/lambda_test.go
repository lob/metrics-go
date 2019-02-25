package lambda

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testMetric = "test_metric"
const testCount = int64(1)
const testValue = float64(50)
const testRate = float64(1)
const testType = "test_type"
const testTag = "foo:bar"

var testTags = []string{testTag}

type mockWriteCloser struct {
	t      *testing.T
	buffer *bytes.Buffer
	closed bool
}

func (w *mockWriteCloser) Write(p []byte) (n int, err error) {
	return w.buffer.Write(p)
}

func (w *mockWriteCloser) Close() error {
	w.closed = true
	return nil
}

func newMockedClient(t *testing.T) *Client {
	w := &mockWriteCloser{t, new(bytes.Buffer), false}
	c, err := New(w)
	assert.NoError(t, err)

	c.Namespace = "testing"
	return c
}

func TestNewLambda(t *testing.T) {
	w := &mockWriteCloser{t, new(bytes.Buffer), false}

	t.Run("create new lambda Client", func(t *testing.T) {
		c, err := New(w)
		assert.NoError(t, err)

		assert.Equal(t, w, c.writer)
	})
}

func TestClose(t *testing.T) {
	t.Run("calls Close function and closes the WriteCloser", func(t *testing.T) {
		mc := newMockedClient(t)

		err := mc.Close()
		assert.NoError(t, err)

		w := mc.writer.(*mockWriteCloser)
		assert.Equal(t, true, w.closed)
	})
}

func TestCount(t *testing.T) {
	t.Run("calls Count function and calls send", func(t *testing.T) {
		mc := newMockedClient(t)

		err := mc.Count(testMetric, testCount, testTags, testRate)
		assert.NoError(t, err)

		w := mc.writer.(*mockWriteCloser)
		got := w.buffer.String()
		assert.Equal(t, strings.Contains(got, "MONITORING"), true)
		assert.Equal(t, strings.Contains(got, "count"), true)
		assert.Equal(t, strings.Contains(got, strconv.FormatInt(testCount, 10)), true)
	})
}

func TestGauge(t *testing.T) {
	t.Run("calls Gauge function and calls send", func(t *testing.T) {
		mc := newMockedClient(t)

		err := mc.Gauge(testMetric, testValue, testTags, testRate)
		assert.NoError(t, err)

		w := mc.writer.(*mockWriteCloser)
		got := w.buffer.String()
		fmt.Println(strconv.FormatFloat(50, 'f', -1, 64))
		assert.Equal(t, strings.Contains(got, "MONITORING"), true)
		assert.Equal(t, strings.Contains(got, "gauge"), true)
		assert.Equal(t, strings.Contains(got, strconv.FormatFloat(testValue, 'f', -1, 64)), true)
	})
}

func TestHistogram(t *testing.T) {
	t.Run("calls Histogram function and calls send", func(t *testing.T) {
		mc := newMockedClient(t)

		err := mc.Histogram(testMetric, testValue, testTags, testRate)
		assert.NoError(t, err)

		w := mc.writer.(*mockWriteCloser)
		got := w.buffer.String()
		fmt.Println(got)
		assert.Equal(t, strings.Contains(got, "MONITORING"), true)
		assert.Equal(t, strings.Contains(got, "histogram"), true)
		assert.Equal(t, strings.Contains(got, strconv.FormatFloat(testValue, 'f', -1, 64)), true)
	})
}

func TestSend(t *testing.T) {
	w := &mockWriteCloser{t, new(bytes.Buffer), false}

	t.Run("calls send function and writes a string in the correct format", func(t *testing.T) {
		c, err := New(w)
		assert.NoError(t, err)

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
			"MONITORING|%d|%s|%s|%s.%s|#%s",
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
