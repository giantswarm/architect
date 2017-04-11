// Package giantswarm provides configuration structures for GiantSwarm services.
package giantswarm

import "github.com/giantswarm/architect/configuration/giantswarm/api"
import "github.com/giantswarm/architect/configuration/giantswarm/passage"
import "github.com/giantswarm/architect/configuration/giantswarm/desmotes"
import "github.com/giantswarm/architect/configuration/giantswarm/happa"

// GiantSwarm holds configuration for GiantSwarm services.
type GiantSwarm struct {
	// API holds configuration for the GiantSwarm API.
	api.API
	// Passage holds configuration for Passage.
	passage.Passage
	// Desmotes holds configuration for Desmotes.
	desmotes.Desmotes
	// Happa holds configuration for Happa.
	happa.Happa
}
