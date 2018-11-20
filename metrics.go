package metrics

import (
	"fmt"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/lob/assets-proxy/pkg/config"
)

type statsdClient interface {
	Histogram(name string, value float64, tags []string, rate float64) error
	Count(name string, value int64, tags []string, rate float64) error
}

// Metrics functions for metrics clients.
type Metrics struct {
	client statsdClient
}

const namespace = "assets_proxy."

// New sets up metric package with a Datadog client.
func New(cfg config.Config) (Metrics, error) {
	address := fmt.Sprintf("%s:%d", cfg.StatsdHost, cfg.StatsdPort)

	client, err := statsd.New(address)
	if err != nil {
		return Metrics{}, err
	}

	client.Namespace = namespace
	client.Tags = []string{
		fmt.Sprintf("environment:%s", cfg.Environment),
		fmt.Sprintf("container:%s", cfg.Hostname),
		fmt.Sprintf("release:%s", cfg.Release),
	}

	return Metrics{client}, nil
}

// Count increments an event counter in Datadog while disregarding potential
// errors.
func (m *Metrics) Count(name string, count int64, tags ...string) {
	m.client.Count(name, count, tags, 1) // nolint:gosec
}

// Histogram sends statistical distribution data to Datadog while disregarding
// potential errors.
func (m *Metrics) Histogram(name string, value float64, tags ...string) {
	m.client.Histogram(name, value, tags, 1) // nolint:gosec
}
