package metrics

import (
	"errors"
	"fmt"
	"strings"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/lob/metrics-go/pkg/lambda"
)

type metricsClient interface {
	Count(name string, count int64, tags []string, rate float64) error
	Gauge(name string, value float64, tags []string, rate float64) error
	Histogram(name string, value float64, tags []string, rate float64) error
	Close() error
}

// StatsReporter provides the ability to report metrics to a client.
type StatsReporter struct {
	client metricsClient
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

	var client metricsClient
	if cfg.Lambda {
		lambda, err := lambda.New(cfg.LambdaLogger)
		if err != nil {
			return &StatsReporter{}, err
		}

		lambda.Namespace = cfg.Namespace
		lambda.Tags = []string{
			fmt.Sprintf("environment:%s", cfg.Environment),
			fmt.Sprintf("container:%s", cfg.Hostname),
			fmt.Sprintf("release:%s", cfg.Release),
		}

		client = lambda
	} else {
		address := fmt.Sprintf("%s:%d", cfg.StatsdHost, cfg.StatsdPort)
		statsd, err := statsd.New(address)
		if err != nil {
			return &StatsReporter{}, err
		}

		statsd.Namespace = cfg.Namespace
		statsd.Tags = []string{
			fmt.Sprintf("environment:%s", cfg.Environment),
			fmt.Sprintf("container:%s", cfg.Hostname),
			fmt.Sprintf("release:%s", cfg.Release),
		}

		client = statsd
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
	m.client.Close() // nolint:gosec
}
