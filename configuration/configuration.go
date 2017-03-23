// Package configuration provides structures for configuring a G8S installation.
// The entire configuration structure is versioned.
// The versioning contract is that fields can be added to a version,
// but not removed or changed within a version.
package configuration

import (
	"github.com/giantswarm/architect/configuration/auth"
	"github.com/giantswarm/architect/configuration/cluster"
	"github.com/giantswarm/architect/configuration/giantswarm"
	"github.com/giantswarm/architect/configuration/monitoring"
)

// Installation holds all the configuration for a G8S installation.
type Installation struct {
	// V1 is the version 1 of the configuration structure.
	V1
}

// V1 is the version 1 of the configuration structure.
type V1 struct {
	// Auth holds configuration for authentication/authorization services.
	auth.Auth

	// GiantSwarm holds configuration for GiantSwarm services.
	giantswarm.GiantSwarm

	// Guest holds configuration for guest clusters.
	cluster.Guest

	// Monitoring holds configuration for monitoring services.
	monitoring.Monitoring
}
