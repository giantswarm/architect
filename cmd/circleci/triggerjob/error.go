package triggerjob

import (
	"github.com/giantswarm/microerror"
)

var missingBranchError = &microerror.Error{
	Kind: "missingBranchError",
}

func IsMissingBranchnError(err error) bool {
	return microerror.Cause(err) == missingBranchError
}

var missingOrganisationError = &microerror.Error{
	Kind: "missingOrganisationError",
}

func IsMissingOrganisationnError(err error) bool {
	return microerror.Cause(err) == missingOrganisationError
}

var missingProjectError = &microerror.Error{
	Kind: "missingProjectError",
}

func IsMissingProjectnError(err error) bool {
	return microerror.Cause(err) == missingProjectError
}
