package chart

import (
	"github.com/giantswarm/microerror"
)

var missingGithubTokenError = &microerror.Error{
	Kind: "missingGithubTokenError",
}

func IsMissingGithubTokenError(err error) bool {
	return microerror.Cause(err) == missingGithubTokenError
}

var createReleaseDirError = &microerror.Error{
	Kind: "createReleaseDirError",
}

func IsCreateReleaseDirError(err error) bool {
	return microerror.Cause(err) == createReleaseDirError
}

var createWorkflowError = &microerror.Error{
	Kind: "createWorkflowError",
}

func IsCreateWorkflowError(err error) bool {
	return microerror.Cause(err) == createWorkflowError
}
