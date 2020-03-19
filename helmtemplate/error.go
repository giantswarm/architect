package helmtemplate

import (
	"github.com/giantswarm/microerror"
)

var executionFailedError = &microerror.Error{
	Kind: "executionFailedError",
}

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var validationFailedError = &microerror.Error{
	Kind: "validationFailedError",
}

// IsValidationFailedError asserts validationFailedError.
func IsValidationFailedError(err error) bool {
	return microerror.Cause(err) == validationFailedError
}
