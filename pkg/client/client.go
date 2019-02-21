package client

import (
	"errors"
	"fmt"
	"strings"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/lob/metrics-go/pkg/config"
	"github.com/lob/metrics-go/pkg/logger"
)

type metricsClient interface {
	Count(name string, count int64, tags []string, rate float64) error
	Gauge(name string, value float64, tags []string, rate float64) error
	Histogram(name string, value float64, tags []string, rate float64) error
	Close() error
}

type Client struct {
	metricsClient
}

// New sets up metric package with a Datadog client.
func New(cfg config.Config) (*Client, error) {
	if cfg.Namespace == "" {
		// Namespace must be populated
		return nil, errors.New("Namespace must be provided")
	}
	if !strings.HasSuffix(cfg.Namespace, ".") {
		cfg.Namespace = fmt.Sprintf("%s.", cfg.Namespace)
	}

	var client *Client
	if cfg.Log {
		logger := logger.New(cfg.Logger)
		logger.Namespace = cfg.Namespace
		logger.Tags = []string{
			fmt.Sprintf("environment:%s", cfg.Environment),
			fmt.Sprintf("release:%s", cfg.Release),
		}

		client = &Client{metricsClient: logger}
	} else {
		address := fmt.Sprintf("%s:%d", cfg.StatsdHost, cfg.StatsdPort)

		statsd, err := statsd.New(address)
		if err != nil {
			return &Client{}, err
		}

		statsd.Namespace = cfg.Namespace
		statsd.Tags = []string{
			fmt.Sprintf("environment:%s", cfg.Environment),
			fmt.Sprintf("container:%s", cfg.Hostname),
			fmt.Sprintf("release:%s", cfg.Release),
		}

		client = &Client{metricsClient: statsd}
	}

	return client, nil
}
