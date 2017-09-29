package template

import (
	"github.com/giantswarm/microerror"
)

var multipleHelmChartsError = microerror.New("multiple helm charts")

// IsMultipleHelmCharts asserts multipleHelmChartsError.
func IsMultipleHelmChart(err error) bool {
	return microerror.Cause(err) == multipleHelmChartsError
}

var incorrectShaError = microerror.New("incorrect sha")

// IsIncorrectSha asserts incorrectShaError.
func IsIncorrectSha(err error) bool {
	return microerror.Cause(err) == incorrectShaError
}

var incorrectValueError = microerror.New("incorrect value")

// IsIncorrectValueError asserts incorrectValueError.
func IsIncorrectValue(err error) bool {
	return microerror.Cause(err) == incorrectValueError
}

var multipleFilesFoundInResourcesError = microerror.New("multiple files found in resources")

// IsMultipleFilesFoundInResources asserts multipleFilesFoundInResourcesError.
func IsMultipleFilesFoundInResources(err error) bool {
	return microerror.Cause(err) == multipleFilesFoundInResourcesError
}

var resourceNotFoundError = microerror.New("resource not found")

// IsResourceNotFound asserts resourceNotFoundError.
func IsResourceNotFound(err error) bool {
	return microerror.Cause(err) == resourceNotFoundError
}

var incorrectTemplatingError = microerror.New("incorrect templating")

// IsIncorrectTemplating asserts incorrectTemplatingError.
func IsIncorrectTemplating(err error) bool {
	return microerror.Cause(err) == incorrectTemplatingError
}
