package appcr

import (
	"strconv"
	"strings"

	"github.com/giantswarm/microerror"
)

func validateConfigVersion(v string) error {
	split := strings.SplitN(v, ".", 2)
	if _, err := strconv.Atoi(split[0]); err != nil {
		// If the string doesn't start with a number and dot assume
		// this is a valid branch name.
		return nil
	}
	if len(split) != 2 || split[1] != "x.x" {
		return microerror.Maskf(executionFailedError, "configuration version starting with a number followed by dot is supposed to end with %q got %q", "x.x", v)
	}

	return nil
}
