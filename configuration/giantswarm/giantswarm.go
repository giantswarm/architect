// Package giantswarm provides configuration structures for GiantSwarm services.
package giantswarm

import "github.com/giantswarm/architect/configuration/giantswarm/api"

// GiantSwarm holds configuration for GiantSwarm services.
type GiantSwarm struct {
	// API holds configuration for the GiantSwarm API.
	api.API
}
