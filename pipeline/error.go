package pipeline

import "github.com/giantswarm/microerror"

var incorrectChartVersionError = &microerror.Error{
	Kind: "incorrectChartVersionError",
}

// IsIncorrectChartVersion asserts incorrectChartVersion.
func IsIncorrectChartVersion(err error) bool {
	return microerror.Cause(err) == incorrectChartVersionError
}
