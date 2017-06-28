// Package auth provides configuration structures for Kubernetes installations.
package kubernetes

// Kubernetes holds configuration for guest cluster's Kubernetes installation.
type Kubernetes struct {
	// API holds configuration for the Kubernetes API
	API

	// IngressController holds configuration for the Ingress Controller
	IngressController
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

// IngressController holds configuration for the Ingress Controller.
type IngressController struct {
	// BaseDomain is the base domain for the Ingress Controller recordset.
	// In some installations a separate domain is needed for security.
	// e.g. API 'api.cluster.k8s.fra-1.giantswarm.io'
	// Ingress ' ingress.cluster.k8s.fra-1.gigantic.io'
	BaseDomain string
}