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
	versionMismatchError = &microerror.Error{
		Kind: "versionMismatchError",
	}
)

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

// IsVersionValidationError asserts versionValidationError.
func IsVersionValidationError(err error) bool {
	return microerror.Cause(err) == versionValidationError
}

// IsVersionMismatchError asserts versionMismatchError.
func IsVersionMismatchError(err error) bool {
	return microerror.Cause(err) == versionMismatchError
}
