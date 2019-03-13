package lambda

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"time"
)

// Client writes metrics in a Lambda format and
// implements the metricsClient interface.
type Client struct {
	writer    io.WriteCloser
	Namespace string
	Tags      []string
}

// New returns a new lambda Client that uses the provided WriteCloser.
func New(w io.WriteCloser) (*Client, error) {
	if w == nil {
		return nil, errors.New("invalid writer")
	}

	return &Client{writer: w}, nil
}

// Close wraps the writers Close method.
func (l *Client) Close() error {
	return l.writer.Close()
}

// Count converts count to a string and sends the metric.
func (l *Client) Count(name string, count int64, tags []string, rate float64) error {
	return l.send(name, strconv.FormatInt(count, 10), tags, "count")
}

// Gauge converts the gauge value to a string and sends the metric.
func (l *Client) Gauge(name string, value float64, tags []string, rate float64) error {
	return l.send(name, strconv.FormatFloat(value, 'f', -1, 64), tags, "gauge")
}

// Histogram converts the histogram value to a string and sends the metric.
func (l *Client) Histogram(name string, value float64, tags []string, rate float64) error {
	return l.send(name, strconv.FormatFloat(value, 'f', -1, 64), tags, "histogram")
}

// nolint:gosec
func (l *Client) send(name string, value string, tags []string, metricType string) error {
	now := time.Now().Unix()

	var buffer bytes.Buffer
	buffer.WriteString("MONITORING|")
	buffer.WriteString(strconv.FormatInt(int64(now), 10))
	buffer.WriteRune('|')
	buffer.WriteString(value)
	buffer.WriteRune('|')
	buffer.WriteString(metricType)
	buffer.WriteRune('|')
	buffer.WriteString(l.Namespace)
	buffer.WriteString(name)
	buffer.WriteRune('|')

	tgs := make([]string, 0, len(tags))
	tgs = append(tgs, l.Tags...)
	tgs = append(tgs, tags...)

	buffer.WriteRune('#')
	buffer.WriteString(tgs[0])
	for _, tag := range tgs[1:] {
		buffer.WriteString(",")
		buffer.WriteString(tag)
	}

	_, err := l.writer.Write(buffer.Bytes())
	return err
}
