// Package configuration provides structures for configuring a G8S installation.
// The entire configuration structure is versioned.
// The versioning contract is that fields can be added to a version,
// but not removed or changed within a version.
package configuration

// Installation holds all the configuration for a G8S installation.
type Installation struct {
	V1
}

// V1 is the version 1 of the configuration structure.
type V1 struct {
	APIServices
	GuestClusters
	Monitoring
	Vault
}
