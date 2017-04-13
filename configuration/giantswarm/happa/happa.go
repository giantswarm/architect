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

	// CreateClusterWorkerType controls the type of form to show when creating a
	// cluster. Valid values are 'aws' or 'kvm'.
	CreateClusterWorkerType string
}
