package utils

import (
	"github.com/juju/errgo"
)

var sourceNotDirectoryError = errgo.New("source not directory")

// IsSourceNotDirectory asserts sourceNotDirectoryError
func IsSourceNotDirectory(err error) bool {
	return errgo.Cause(err) == sourceNotDirectoryError
}

var destinationExistsError = errgo.New("destination exists")

// IsDestinationExists asserts destinationExistsError.
func IsDestinationExists(err error) bool {
	return errgo.Cause(err) == destinationExistsError
}
