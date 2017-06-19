package ingress

type Ingress struct {
	PortRange PortRange
}

type PortRange struct {
	Min, Max int
}
