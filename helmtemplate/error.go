package helmtemplate

import (
	"github.com/giantswarm/microerror"
)

var (
	executionFailedError = &microerror.Error{
		Kind: "executionFailedError",
	}

	invalidConfigError = &microerror.Error{
		Kind: "invalidConfigError",
	}
)

// IsExecutionFailedError asserts executionFailedError.
func IsExecutionFailedError(err error) bool {
	return microerror.Cause(err) == executionFailedError
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}
