package calico // Package calico provides configuration structures for Calico.

// Calico holds configuration for guest cluster calico settings.
type Calico struct {
	// Subnet for calico, eg. 192.168.0.0.
	Subnet string
	// Cidr for calico subnets, eg. 16.
	CIDR string
}
