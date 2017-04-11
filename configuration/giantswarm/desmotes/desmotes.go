// Package desmotes provides configuration structures for a Desmotes service.
package desmotes

import (
	"net/url"
)

// Desmotes holds configuration for a Desmotes service.
type Desmotes struct {
	// Address is the URL to Desmotes.
	// e.g: 'https://desmotes-g8s.giantswarm.io'
	Address url.URL
}
