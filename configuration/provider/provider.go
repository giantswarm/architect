// Package provider exposes provider specific configurations used to express
// differences in the environmental virtualization technologies.
package provider

// kind is a private type to ensure only providers defined in this package can
// be applied to installation configurations. That prevents other packages
// screwing around with provider configurations.
type kind string

const (
	// AWS represents a G8S installation being deployed on AWS where EC2 provides
	// the virtualization.
	AWS kind = "aws"
	// KVM represents a G8S installation being deployed on bare metal where KVM
	// provides the virtualization.
	KVM kind = "kvm"
)

// Provider holds configuration for monitoring services.
type Provider struct {
	// Kind is the provider kind.
	Kind kind
}
