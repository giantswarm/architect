package pipeline

import "github.com/giantswarm/microerror"

var incorrectChartVersionError = microerror.New("incorrect chart version")

// IsIncorrectChartVersion asserts incorrectChartVersion.
func IsIncorrectChartVersion(err error) bool {
	return microerror.Cause(err) == incorrectChartVersionError
}
