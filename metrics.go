package metrics

import (
	"github.com/lob/metrics-go/pkg/client"
	"github.com/lob/metrics-go/pkg/config"
)

// StatsReporter provides the ability to report metrics to statsd
type StatsReporter struct {
	client *client.Client
}

// New sets up metric package with a Datadog client.
func New(cfg config.Config) (*StatsReporter, error) {
	client, err := client.New(cfg)
	if err != nil {
		return &StatsReporter{}, err
	}

	return &StatsReporter{client}, nil
}

// Count increments an event counter in Datadog while disregarding potential
// errors.
func (m *StatsReporter) Count(name string, count int64, tags ...string) {
	m.client.Count(name, count, tags, 1) // nolint:gosec
}

// Gauge sends a current value for an event counter in Datadog
// while disregarding potential errors.
func (m *StatsReporter) Gauge(name string, value float64, tags ...string) {
	m.client.Gauge(name, value, tags, 1) // nolint:gosec
}

// Histogram sends statistical distribution data to Datadog while disregarding
// potential errors.
func (m *StatsReporter) Histogram(name string, value float64, tags ...string) {
	m.client.Histogram(name, value, tags, 1) // nolint:gosec
}

// Close closes the metrics client while disregarding potential errors.
func (m *StatsReporter) Close() {
	m.client.Close()
}
