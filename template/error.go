package template

import (
	"github.com/giantswarm/microerror"
)

var incorrectShaError = &microerror.Error{
	Kind: "incorrectShaError",
}

// IsIncorrectSha asserts incorrectShaError.
func IsIncorrectSha(err error) bool {
	return microerror.Cause(err) == incorrectShaError
}

var incorrectValueError = &microerror.Error{
	Kind: "incorrectValueError",
}

// IsIncorrectValue asserts incorrectValueError.
func IsIncorrectValue(err error) bool {
	return microerror.Cause(err) == incorrectValueError
}

var multipleFilesFoundInResourcesError = &microerror.Error{
	Kind: "multipleFilesFoundInResourcesError",
}

// IsMultipleFilesFoundInResources asserts multipleFilesFoundInResourcesError.
func IsMultipleFilesFoundInResources(err error) bool {
	return microerror.Cause(err) == multipleFilesFoundInResourcesError
}

var resourceNotFoundError = &microerror.Error{
	Kind: "resourceNotFoundError",
}

// IsResourceNotFound asserts resourceNotFoundError.
func IsResourceNotFound(err error) bool {
	return microerror.Cause(err) == resourceNotFoundError
}

var incorrectTemplatingError = &microerror.Error{
	Kind: "incorrectTemplatingError",
}

// IsIncorrectTemplating asserts incorrectTemplatingError.
func IsIncorrectTemplating(err error) bool {
	return microerror.Cause(err) == incorrectTemplatingError
}
