package calico // Package calico provides configuration structures for Calico.

// Calico holds configuration for guest cluster calico settings.
type Calico struct {
	// subnet for calico
	// ie 192.168.0.0
	Subnet string
	// cidr for calico sibnets
	// ie 16
	CIDR string
}
