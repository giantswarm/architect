package argoapp

import "github.com/giantswarm/microerror"

// executionFailedError should never be matched against and therefore there is
// no matcher implement. For further information see:
//
//     https://github.com/giantswarm/fmt/blob/master/go/errors.md#matching-errors
//
var executionFailedError = &microerror.Error{
	Kind: "executionFailedError",
}
