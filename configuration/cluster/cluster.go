// Package cluster provides configuration structures for guest clusters.
package cluster

// Guest holds the configuration for the installation's guest clusters.
type Guest struct {
	// Hyperkube holds configuration for the guest cluster's
	// Hyperkube settings.
	Hyperkube

	// Kubernetes holds configuration for the guest cluster's
	// Kubernetes installation.
	Kubernetes
}

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

// Hyperkube holds configuration for Hyperkube settings.
type Hyperkube struct {
	// Version is the version of Hyperkube.
	// e.g: 'v1.5.2_coreos.0'
	Version string
}
