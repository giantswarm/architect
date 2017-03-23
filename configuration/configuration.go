// Package configuration provides structures for configuring a G8S installation.
// The entire configuration structure is versioned.
// The versioning contract is that fields can be added to a version,
// but not removed or changed within a version.
package configuration

import (
	"github.com/giantswarm/architect/configuration/apiservices"
	"github.com/giantswarm/architect/configuration/guestclusters"
	"github.com/giantswarm/architect/configuration/monitoring"
	"github.com/giantswarm/architect/configuration/vault"
)

// Installation holds all the configuration for a G8S installation.
type Installation struct {
	V1
}

// V1 is the version 1 of the configuration structure.
type V1 struct {
	apiservices.APIServices
	guestclusters.GuestClusters
	monitoring.Monitoring
	vault.Vault
}
