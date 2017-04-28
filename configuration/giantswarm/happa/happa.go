// Happa package provides configuration structures for a Happa service
package happa

import (
	"net/url"
)

// Happa holds configuration for a Happa service.
type Happa struct {
	// Address is the URL to Happa.
	// e.g: 'https://happa-g8s.giantswarm.io'
	Address url.URL
}
