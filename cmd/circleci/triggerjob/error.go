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

var missingRepositoryError = &microerror.Error{
	Kind: "missingRepositoryError",
}

func IsMissingRepositorynError(err error) bool {
	return microerror.Cause(err) == missingRepositoryError
}
