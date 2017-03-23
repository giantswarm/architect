// Package apiservices provides configuration structures for API services.
package apiservices

import (
	"net/url"
)

// APIServices holds configuration for API services.
type APIServices struct {
	GSAPI
}

// GSAPI holds configuration for the installation's Giant Swarm API.
type GSAPI struct {
	// Address is the URL of the Giant Swarm API.
	// e.g: 'https://api-g8s.giantswarm.io'
	Address url.URL
}
