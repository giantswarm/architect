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

var emptyKubernetesResourcesDirectoryPath = errgo.New("empty kubernetes resources directory path")

// IsEmptyKubernetesResourcesDirectoryPath asserts emptyKubernetesResourcesDirectoryPath
func IsEmptyKubernetesResourcesDirectoryPath(err error) bool {
	return errgo.Cause(err) == emptyKubernetesResourcesDirectoryPath
}

var nilTemplateStructError = errgo.New("nil template struct")

// IsNilTemplateStruct asserts nilTemplateStruct.
func IsNilTemplateStruct(err error) bool {
	return errgo.Cause(err) == nilTemplateStructError
}

var notStringTypeError = errgo.New("not string type")

// IsNotStringType asserts notStringTypeError.
func IsNotStringType(err error) bool {
	return errgo.Cause(err) == notStringTypeError
}
