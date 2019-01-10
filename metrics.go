package metrics

import (
	"errors"
	"fmt"
	"strings"

	"github.com/DataDog/datadog-go/statsd"
)

type statsdClient interface {
	Histogram(name string, value float64, tags []string, rate float64) error
	Count(name string, value int64, tags []string, rate float64) error
}

// StatsReporter provides the ability to report metrics to statsd
type StatsReporter struct {
	client statsdClient
}

// New sets up metric package with a Datadog client.
func New(cfg Config) (*StatsReporter, error) {
	if cfg.Namespace == "" {
		// Namespace must be populated
		return nil, errors.New("Namespace must be provided")
	}
	if !strings.HasSuffix(cfg.Namespace, ".") {
		cfg.Namespace = fmt.Sprintf("%s.", cfg.Namespace)
	}

	address := fmt.Sprintf("%s:%d", cfg.StatsdHost, cfg.StatsdPort)

	client, err := statsd.New(address)
	if err != nil {
		return &StatsReporter{}, err
	}

	client.Namespace = cfg.Namespace
	client.Tags = []string{
		fmt.Sprintf("environment:%s", cfg.Environment),
		fmt.Sprintf("container:%s", cfg.Hostname),
		fmt.Sprintf("release:%s", cfg.Release),
	}

	return &StatsReporter{client}, nil
}

// Count increments an event counter in Datadog while disregarding potential
// errors.
func (m *StatsReporter) Count(name string, count int64, tags ...string) {
	m.client.Count(name, count, tags, 1) // nolint:gosec
}

// Histogram sends statistical distribution data to Datadog while disregarding
// potential errors.
func (m *StatsReporter) Histogram(name string, value float64, tags ...string) {
	m.client.Histogram(name, value, tags, 1) // nolint:gosec
}
