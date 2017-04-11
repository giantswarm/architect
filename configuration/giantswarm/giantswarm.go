// Package giantswarm provides configuration structures for GiantSwarm services.
package giantswarm

import "github.com/giantswarm/architect/configuration/giantswarm/api"
import "github.com/giantswarm/architect/configuration/giantswarm/passage"
import "github.com/giantswarm/architect/configuration/giantswarm/desmotes"

// GiantSwarm holds configuration for GiantSwarm services.
type GiantSwarm struct {
	// API holds configuration for the GiantSwarm API.
	api.API
	passage.Passage
	desmotes.Desmotes
}
