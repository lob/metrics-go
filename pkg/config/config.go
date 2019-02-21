package config

import "io"

// Config stores configuration necessary for connecting to statsd and reporting metrics.
type Config struct {
	Environment string
	Hostname    string
	Namespace   string
	Release     string
	StatsdHost  string
	StatsdPort  int
	Log         bool
	Logger      io.WriteCloser
}
