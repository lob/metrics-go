package metrics

import (
	"time"
)

// Timer facilitates timing and tagging a metric before sending it to Datadog
// as a Histogram datapoint.
type Timer struct {
	name    string
	metrics *Metrics
	begin   time.Time
	tags    []string
}

// NewTimer returns a Timer object with a set start time
func (m *Metrics) NewTimer(name string, tags ...string) Timer {
	return Timer{
		begin:   time.Now(),
		metrics: m,
		name:    name,
		tags:    tags,
	}
}

// End ends a Timer and sends the metric and duration to Datadog as a
// Histogram datapoint.
func (t *Timer) End(additionalTags ...string) float64 {
	duration := time.Since(t.begin)
	durationInMS := float64(duration / time.Millisecond)

	t.tags = append(t.tags, additionalTags...)

	t.metrics.Histogram(t.name, durationInMS, t.tags...)

	return durationInMS
}
