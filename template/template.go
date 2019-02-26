// Package template provides functions for templating Kubernetes resources,
// with standard G8S configuration.
package template

import (
	"os"
)

var (
	// permission is a default permission to use for templated files
	permission os.FileMode = 0644
)

// BuildInfo holds information concerning the current build.
type BuildInfo struct {
	// DockerTag is the docker image tag.
	DockerTag string
	// SHA is the SHA-1 tag of the commit we are building for.
	SHA string
	// Version is the version of the commit being built.
	Version string
}

// TemplateConfiguration holds both build info, and configuration info.
type TemplateConfiguration struct {
	// BuildInfo is the configuration for the current build
	BuildInfo
}
