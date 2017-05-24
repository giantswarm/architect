package utils

import (
	"github.com/juju/errgo"
)

var sourceIsNotDirectoryError = errgo.New("source is not a directory")

// IsSourceIsNotDirectory asserts sourceIsNotDirectoryError.
func IsSourceIsNotDirectory(err error) bool {
	return errgo.Cause(err) == sourceIsNotDirectoryError
}

var destinationExistsError = errgo.New("destination already exists")

// IsDestinationExists asserts destinationExistsError.
func IsDestinationExists(err error) bool {
	return errgo.Cause(err) == destinationExistsError
}
