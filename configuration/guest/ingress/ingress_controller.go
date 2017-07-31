// Package ingress provides configuration structures for the Ingress Controller.
package ingress

// kind is a private type to ensure only versions defined in this package can
// be applied to installation configurations. That prevents other packages
// screwing around with version configurations.
type kind string

const (
	Version kind = "0.9.0-beta.11"
)

// IngressController holds configuration for Ingress Controller settings.
type IngressController struct {
	// Version is the version of Ingress Controller, e.g. '0.9.0-beta.11'.
	Version kind
}
