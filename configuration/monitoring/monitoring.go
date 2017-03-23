// Package monitoring provides configuration structures for monitoring services.
package monitoring

import (
	"net/url"
	"time"
)

// Monitoring holds configuration for monitoring services.
type Monitoring struct {
	Prometheus
	Testbot
}

// Prometheus holds the configuration for the installation's Prometheus setup.
type Prometheus struct {
	// Address is the URL of the Prometheus API.
	// e.g: 'https://prometheus-g8s.giantswarm.io'
	Address url.URL

	// RetentionPeriod is how long to keep Prometheus data for.
	// e.g: '336h'
	RetentionPeriod time.Duration
}

// Testbot holds the configuration for the installation's Testbot setup.
type Testbot struct {
	// Interval is the time between testbot runs.
	// e.g: '5m'
	Interval time.Duration
}
