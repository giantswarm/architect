// Package auth provides configuration structures for Kubernetes installations.
package kubernetes

// Kubernetes holds configuration for guest cluster's Kubernetes installation.
type Kubernetes struct {
	// API holds configuration for the Kubernetes API
	API
}

// API holds configuration for the Kubernetes API
type API struct {
	// EndpointBase is the base URL used to configure guest clusters.
	// e.g: 'g8s.fra-1.giantswarm.io'
	// For example, used in:
	// - Guest cluster endpoint format: 'https://api.%s.g8s.fra-1.giantswarm.io'
	// - Common name format: '%s.g8s.fra-1.giantswarm.io'
	EndpointBase string
}
