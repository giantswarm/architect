package configuration

// GuestClusters holds the configuration for the installation's
// guest clusters.
type GuestClusters struct {
	// APIEndpointBase is the base URL used to configure guest clusters.
	// e.g: 'g8s.fra-1.giantswarm.io'
	// For example, used in:
	// - Guest cluster endpoint format: 'https://api.%s.g8s.fra-1.giantswarm.io'
	// - Common name format: '%s.g8s.fra-1.giantswarm.io'
	APIEndpointBase string

	// HyperkubeVersion is the version of hyperkube images to use.
	// e.g: 'v1.5.2_coreos.0'
	HyperkubeVersion string
}
