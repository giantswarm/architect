// Package hyperkube provides configuration structures for Hyperkube.
package hyperkube

// kind is a private type to ensure only versions defined in this package can
// be applied to installation configurations. That prevents other packages
// screwing around with version configurations.
type kind string

const (
	Version kind = "v1.6.4_coreos.0"
)

// Hyperkube holds configuration for Hyperkube settings.
type Hyperkube struct {
	// Version is the version of Hyperkube, e.g. 'v1.6.4_coreos.0'.
	Version kind
}
