package metrics

import (
	"fmt"

	"github.com/DataDog/datadog-go/statsd"
)

type statsdClient interface {
	Histogram(name string, value float64, tags []string, rate float64) error
	Count(name string, value int64, tags []string, rate float64) error
}

type metrics struct {
	client statsdClient
}

// Metrics defines the interface for metrics producers
type Metrics interface {
	Count(name string, count int64, tags ...string)
	Histogram(name string, value float64, tags ...string)
	NewTimer(name string, tags ...string) Timer
}

// New sets up metric package with a Datadog client.
func New(cfg Config) (Metrics, error) {
	address := fmt.Sprintf("%s:%d", cfg.StatsdHost, cfg.StatsdPort)

	client, err := statsd.New(address)
	if err != nil {
		return &metrics{}, err
	}

	client.Namespace = cfg.Namespace
	client.Tags = []string{
		fmt.Sprintf("environment:%s", cfg.Environment),
		fmt.Sprintf("container:%s", cfg.Hostname),
		fmt.Sprintf("release:%s", cfg.Release),
	}

	return &metrics{client}, nil
}

// Count increments an event counter in Datadog while disregarding potential
// errors.
func (m *metrics) Count(name string, count int64, tags ...string) {
	m.client.Count(name, count, tags, 1) // nolint:gosec
}

// Histogram sends statistical distribution data to Datadog while disregarding
// potential errors.
func (m *metrics) Histogram(name string, value float64, tags ...string) {
	m.client.Histogram(name, value, tags, 1) // nolint:gosec
}
