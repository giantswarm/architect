package internal

import (
	"github.com/giantswarm/microerror"
)

var executionFailedError = &microerror.Error{
	Kind: "executionFailedError",
}

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var fileNotFoundError = &microerror.Error{
	Kind: "fileNotFoundError",
}

// IsFileNotFound asserts fileNotFoundError.
func IsFileNotFound(err error) bool {
	return microerror.Cause(err) == fileNotFoundError
}

var nonCanonicalHeadingError = &microerror.Error{
	Kind: "nonCanonicalHeadingError",
}

// IsNonCanonicalHeading asserts nonCanonicalHeadingError.
func IsNonCanonicalHeading(err error) bool {
	return microerror.Cause(err) == nonCanonicalHeadingError
}

var missingStableSectionError = &microerror.Error{
	Kind: "missingStableSectionError",
}

// IsMissingStableSection asserts missingStableSectionError.
func IsMissingStableSection(err error) bool {
	return microerror.Cause(err) == missingStableSectionError
}

var orphanContentError = &microerror.Error{
	Kind: "orphanContentError",
}

// IsOrphanContent asserts orphanContentError.
func IsOrphanContent(err error) bool {
	return microerror.Cause(err) == orphanContentError
}
