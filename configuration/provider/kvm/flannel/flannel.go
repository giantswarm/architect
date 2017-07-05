package flannel

type Flannel struct {
	VNIRange       Range
	NetworkFormat  string
	Interface      string
	PrivateNetwork string
}

type Range struct {
	Min, Max int
}
