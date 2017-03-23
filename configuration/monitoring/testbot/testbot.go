// Package testbot provides configuration structures for the Testbot service.
package testbot

import (
	"time"
)

// Testbot holds the configuration for the installation's Testbot setup.
type Testbot struct {
	// Interval is the time between testbot runs.
	// e.g: '5 * time.Minute'
	Interval time.Duration
}
