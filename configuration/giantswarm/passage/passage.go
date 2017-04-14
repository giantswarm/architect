// Package passage provides configuration structures for a Passage service
package passage

import (
	"net/url"
)

// Passage holds configuration for the a Passage service.
type Passage struct {
	// Address is the URL to Passage.
	// e.g: 'https://passage-g8s.giantswarm.io'
	Address url.URL
}
