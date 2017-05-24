package template

import (
	"github.com/juju/errgo"
)

var multipleHelmChartsError = errgo.New("multiple charts found in helm directory")

// IsMultipleHelmCharts asserts multipleHelmChartsError.
func IsMultipleHelmChart(err error) bool {
	return errgo.Cause(err) == multipleHelmChartsError
}

var incorrectShaError = errgo.New("correct sha not found in chart")

// IsIncorrectSha asserts incorrectShaError.
func IsIncorrectSha(err error) bool {
	return errgo.Cause(err) == incorrectShaError
}

var incorrectValueError = errgo.New("correct value not found in chart")

// IsIncorrectValueError asserts incorrectValueError.
func IsIncorrectValue(err error) bool {
	return errgo.Cause(err) == incorrectValueError
}

var multipleFilesFoundInResourcesError = errgo.New("multiple files found in resources directory")

// IsMultipleFilesFoundInResources asserts multipleFilesFoundInResourcesError.
func IsMultipleFilesFoundInResources(err error) bool {
	return errgo.Cause(err) == multipleFilesFoundInResourcesError
}

var resourceNotFoundError = errgo.New("did not find the required resource")

// IsResourceNotFound asserts resourceNotFoundError.
func IsResourceNotFound(err error) bool {
	return errgo.Cause(err) == resourceNotFoundError
}

var incorrectTemplatingError = errgo.New("did not find required field from template")

// IsIncorrectTemplating asserts incorrectTemplatingError.
func IsIncorrectTemplating(err error) bool {
	return errgo.Cause(err) == incorrectTemplatingError
}
