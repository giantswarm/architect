// Package apiservices provides configuration structures for API services.
package apiservices

import (
	"net/url"
)

// APIServices holds configuration for API services.
type APIServices struct {
	// GSAPI holds configuration for the GiantSwarm API.
	GSAPI
}

// GSAPI holds configuration for the installation's GiantSwarm API.
type GSAPI struct {
	// Address is the URL of the Giant Swarm API.
	// e.g: 'https://api-g8s.giantswarm.io'
	Address url.URL
}
