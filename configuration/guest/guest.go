// Package guest provides configuration structures for guest clusters.
package guest

import (
	"github.com/giantswarm/architect/configuration/guest/calico"
	"github.com/giantswarm/architect/configuration/guest/hyperkube"
	"github.com/giantswarm/architect/configuration/guest/kubectl"
	"github.com/giantswarm/architect/configuration/guest/kubernetes"
)

// Guest holds the configuration for the installation's guest clusters.
type Guest struct {
	// Hyperkube holds configuration for the guest guest cluster's Hyperkube
	// settings.
	hyperkube.Hyperkube

	// Kubectl holds configuration for the guest guest cluster's Kubectl
	// settings.
	kubectl.Kubectl

	// Kubernetes holds configuration for the guest guest cluster's Kubernetes
	// installation.
	kubernetes.Kubernetes

	// Calico holds configuration for calico in guest cluster.
	calico.Calico
}
