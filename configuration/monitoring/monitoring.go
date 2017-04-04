// Package monitoring provides configuration structures for monitoring services.
package monitoring

import (
	"github.com/giantswarm/architect/configuration/monitoring/prometheus"
	"github.com/giantswarm/architect/configuration/monitoring/testbot"
)

// Monitoring holds configuration for monitoring services.
type Monitoring struct {
	// Prometheus holds the configuration for the Prometheus setup.
	prometheus.Prometheus

	// Testbot holds configuration for the Testbot setup.
	testbot.Testbot
}
