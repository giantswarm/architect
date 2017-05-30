// Package kubectl provides configuration structures for Kubectl.
package kubectl

// kind is a private type to ensure only versions defined in this package can
// be applied to installation configurations. That prevents other packages
// screwing around with version configurations.
type kind string

const (
	// This image is for kubectl v1.6.4.
	Version kind = "f51f93c30d27927d2b33122994c0929b3e6f2432"
)

// Kubectl holds configuration for Kubectl settings.
type Kubectl struct {
	// Version is the version of Kubectl, e.g. '1afa480ffb2912fe3605e84fd392c5bd1c9f48b9'.
	Version kind
}
