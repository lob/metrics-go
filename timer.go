package metrics

import (
	"time"
)

// Timer facilitates timing and tagging a metric before sending it to Datadog
// as a Histogram datapoint.
type Timer interface {
	End(...string) float64
}

// timer provides a statsd implementation of the Timer interface
type timer struct {
	name    string
	metrics *metrics
	begin   time.Time
	tags    []string
}

// NewTimer returns a Timer object with a set start time
func (m *metrics) NewTimer(name string, tags ...string) Timer {
	return &timer{
		begin:   time.Now(),
		metrics: m,
		name:    name,
		tags:    tags,
	}
}

// End ends a Timer and sends the metric and duration to Datadog as a
// Histogram datapoint.
func (t *timer) End(additionalTags ...string) float64 {
	duration := time.Since(t.begin)
	durationInMS := float64(duration / time.Millisecond)

	t.tags = append(t.tags, additionalTags...)

	t.metrics.Histogram(t.name, durationInMS, t.tags...)

	return durationInMS
}
