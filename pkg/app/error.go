package app

import (
	"github.com/giantswarm/microerror"
)

var executionFailedError = &microerror.Error{
	Kind: "executionFailedError",
}

func IsExecutionFailedError(err error) bool {
	return microerror.Cause(err) == executionFailedError
}
