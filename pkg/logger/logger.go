package logger

import (
	"bytes"
	"io"
	"strconv"
	"time"
)

type Client struct {
	writer    io.WriteCloser
	Namespace string
	Tags      []string
}

func New(w io.WriteCloser) *Client {
	return &Client{writer: w}
}

func (l *Client) Close() error {
	return l.writer.Close()
}

func (l *Client) Count(name string, count int64, tags []string, rate float64) error {
	return l.send(name, strconv.FormatInt(count, 10), tags, "count")
}

func (l *Client) Gauge(name string, value float64, tags []string, rate float64) error {
	return l.send(name, strconv.FormatFloat(value, 'f', -1, 64), tags, "gauge")
}

func (l *Client) Histogram(name string, value float64, tags []string, rate float64) error {
	return l.send(name, strconv.FormatFloat(value, 'f', -1, 64), tags, "histogram")
}

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
	buffer.WriteRune('.')
	buffer.WriteString(name)
	buffer.WriteRune('|')

	tgs := make([]string, 0, len(tags))
	tgs = append(tgs, l.Tags...)
	tgs = append(tgs, tags...)

	buffer.WriteString("|#")
	buffer.WriteString(tgs[0])
	for _, tag := range tgs[1:] {
		buffer.WriteString(",")
		buffer.WriteString(tag)
	}

	_, err := l.writer.Write(buffer.Bytes())
	return err
}
