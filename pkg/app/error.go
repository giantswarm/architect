package app

import (
	"github.com/giantswarm/microerror"
)

var wrongFormatError = &microerror.Error{
	Kind: "wrongFormatError",
}

func IsWrongFormatError(err error) bool {
	return microerror.Cause(err) == wrongFormatError
}
