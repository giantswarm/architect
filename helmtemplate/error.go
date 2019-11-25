package helmtemplate

import (
	"github.com/giantswarm/microerror"
)

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts incorrectValueError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}
