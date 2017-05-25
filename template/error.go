package template

import (
	"github.com/juju/errgo"
)

var multipleHelmChartsError = errgo.New("multiple helm charts")

// IsMultipleHelmCharts asserts multipleHelmChartsError.
func IsMultipleHelmChart(err error) bool {
	return errgo.Cause(err) == multipleHelmChartsError
}

var incorrectShaError = errgo.New("incorrect sha")

// IsIncorrectSha asserts incorrectShaError.
func IsIncorrectSha(err error) bool {
	return errgo.Cause(err) == incorrectShaError
}

var incorrectValueError = errgo.New("incorrect value")

// IsIncorrectValueError asserts incorrectValueError.
func IsIncorrectValue(err error) bool {
	return errgo.Cause(err) == incorrectValueError
}

var multipleFilesFoundInResourcesError = errgo.New("multiple files found in resources")

// IsMultipleFilesFoundInResources asserts multipleFilesFoundInResourcesError.
func IsMultipleFilesFoundInResources(err error) bool {
	return errgo.Cause(err) == multipleFilesFoundInResourcesError
}

var resourceNotFoundError = errgo.New("resource not found")

// IsResourceNotFound asserts resourceNotFoundError.
func IsResourceNotFound(err error) bool {
	return errgo.Cause(err) == resourceNotFoundError
}

var incorrectTemplatingError = errgo.New("incorrect templating")

// IsIncorrectTemplating asserts incorrectTemplatingError.
func IsIncorrectTemplating(err error) bool {
	return errgo.Cause(err) == incorrectTemplatingError
}
