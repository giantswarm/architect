package flannel

type Flannel struct {
	VNIRange Range
}

type Range struct {
	Min, Max int
}
