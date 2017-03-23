// Package giantswarm provides configuration structures for giantswarm services.
package giantswarm

import (
	"net/url"
)

// GiantSwarm holds configuration for GiantSwarm services.
type GiantSwarm struct {
	// API holds configuration for the GiantSwarm API.
	API
}

// API holds configuration for the installation's GiantSwarm API.
type API struct {
	// Address is the URL of the Giant Swarm API.
	// e.g: 'https://api-g8s.giantswarm.io'
	Address url.URL
}
