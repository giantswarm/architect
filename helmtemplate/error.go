package helmtemplate

import (
	"github.com/giantswarm/microerror"
)

var (
	invalidConfigError = &microerror.Error{
		Kind: "invalidConfigError",
	}
	versionValidationError = &microerror.Error{
		Kind: "versionValidationError",
	}
)

// IsInvalidConfig asserts invalidValueError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

// IsVersionValidationError asserts versionValidationError
func IsVersionValidationError(err error) bool {
	return microerror.Cause(err) == versionValidationError
}
