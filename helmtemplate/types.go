// Package helmtemplate provides functions for templating helm charts
package helmtemplate

import (
	"os"
)

var (
	// permission is a default permission to use for templated files
	permission os.FileMode = 0644
)

// BuildInfo holds information concerning the current build.
type BuildInfo struct {
	// Branch is the name of the branch we're building.
	Branch string
	// SHA is the SHA-1 tag of the commit we are building for.
	SHA string
	// Version is the version of the commit being built.
	Version string
}
