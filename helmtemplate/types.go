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
	// AppVersion is the version read from pkg/project/project.go if it
	// exists or set to the same value as Version otherwise.
	AppVersion string
}

// renderedChart is used for chart validation after it has been filled with
// values
type renderedChart struct {
	Version    string `json:"version,omitempty"`
	AppVersion string `json:"appVersion,omitempty"`
}
