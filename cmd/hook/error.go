package hook

import (
	"github.com/giantswarm/microerror"
)

var gitNoSHAError = &microerror.Error{
	Kind: "gitNoSHAError",
}

func IsGitNoSHAError(err error) bool {
	return microerror.Cause(err) == gitNoSHAError
}

var gitNoBranchError = &microerror.Error{
	Kind: "gitNoBranchError",
}

func IsGitNoBranchError(err error) bool {
	return microerror.Cause(err) == gitNoBranchError
}
