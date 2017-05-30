// Package kubectl provides configuration structures for Kubectl.
package kubectl

// kind is a private type to ensure only versions defined in this package can
// be applied to installation configurations. That prevents other packages
// screwing around with version configurations.
type kind string

const (
	Version kind = "1afa480ffb2912fe3605e84fd392c5bd1c9f48b9"
)

// Kubectl holds configuration for Kubectl settings.
type Kubectl struct {
	// Version is the version of Kubectl, e.g. '1afa480ffb2912fe3605e84fd392c5bd1c9f48b9'.
	Version kind
}
