package release

import (
	"net/http"

	"github.com/giantswarm/microerror"
	"github.com/google/go-github/github"
)

var notFoundError = &microerror.Error{
	Kind: "notFoundError",
}

func IsNotFoundError(err error) bool {
	if microerror.Cause(err) == notFoundError {
		return true
	}

	gErr, ok := err.(*github.ErrorResponse)
	if !ok {
		return false
	}

	if gErr.Response != nil {
		return gErr.Response.StatusCode == http.StatusNotFound
	}

	return false
}
