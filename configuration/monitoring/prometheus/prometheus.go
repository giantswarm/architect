// Package prometheus provides configuration structures for Prometheus setup.
package prometheus

import (
	"net/url"
	"time"
)

// Prometheus holds the configuration for the installation's Prometheus setup.
type Prometheus struct {
	// Address is the URL of the Prometheus API.
	// e.g: 'https://prometheus-g8s.giantswarm.io'
	Address url.URL

	// RetentionPeriod is how long to keep Prometheus data for.
	// e.g: '2 * 7 * 24 * time.Hour'
	RetentionPeriod time.Duration
}
