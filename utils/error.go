package utils

import (
	"github.com/giantswarm/microerror"
)

var sourceNotDirectoryError = microerror.New("source not directory")

// IsSourceNotDirectory asserts sourceNotDirectoryError
func IsSourceNotDirectory(err error) bool {
	return microerror.Cause(err) == sourceNotDirectoryError
}

var destinationExistsError = microerror.New("destination exists")

// IsDestinationExists asserts destinationExistsError.
func IsDestinationExists(err error) bool {
	return microerror.Cause(err) == destinationExistsError
}
