package configuration

// provider is a private type to ensure only providers defined in this package
// can be applied to installation configurations. That prevents other packages
// screwing around with provider configurations.
type provider string

const (
	// ProviderAWS represents a G8S installation being deployed on AWS where EC2
	// provides the virtualization.
	ProviderAWS provider = "aws"
	// ProviderAWS represents a G8S installation being deployed on bare metal
	// where KVM provides the virtualization.
	ProviderKVM provider = "kvm"
)
