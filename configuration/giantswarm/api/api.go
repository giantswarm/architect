// Package api provides configuration structures for giantswarm's API.
package api

import (
	"net/url"
)

// API holds configuration for the installation's GiantSwarm API.
type API struct {
	// Address is the URL of the Giant Swarm API.
	// e.g: 'https://api-g8s.giantswarm.io'
	Address url.URL
}
