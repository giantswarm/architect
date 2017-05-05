package instance

// Allowed returns a list containing all instance types as provided to the
// invocation of Allowed. Using this medthod ensures that only valid instance
// defined by the constants of this package will be used to configure the list
// of allowed instance types of a customer's G8S installation.
func Allowed(instanceTypes ...kind) []kind {
	if len(instanceTypes) == 0 {
		panic("instanceTypes must not be empty")
	}

	return append([]kind{}, instanceTypes...)
}
