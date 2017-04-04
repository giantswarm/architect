// Package cluster provides configuration structures for clusters.
package cluster

import (
	"github.com/giantswarm/architect/configuration/cluster/hyperkube"
	"github.com/giantswarm/architect/configuration/cluster/kubernetes"
)

// Guest holds the configuration for the installation's guest clusters.
type Guest struct {
	// Hyperkube holds configuration for the guest cluster's
	// Hyperkube settings.
	hyperkube.Hyperkube

	// Kubernetes holds configuration for the guest cluster's
	// Kubernetes installation.
	kubernetes.Kubernetes
}
