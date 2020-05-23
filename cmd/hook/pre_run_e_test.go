package hook

import (
	"testing"
)

const (
	FileContentsWIP = `package project

var (
	description string = "The azure-operator manages Kubernetes clusters on Azure."
	gitSHA             = "n/a"
	name        string = "azure-operator"
	source      string = "https://github.com/giantswarm/azure-operator"
	version            = "4.1.0-dev"
)

func Description() string {
	return description
}

func GitSHA() string {
	return gitSHA
}

func Name() string {
	return name
}

func Source() string {
	return source
}

func Version() string {
	return version
}
`
	FileContentsRelease = `package project

import (
	"fmt"
)

var (
	description = "The azure-operator manages Kubernetes clusters on Azure."
	gitSHA      = "n/a"
	name        = "azure-operator"
	source      = "https://github.com/giantswarm/azure-operator"
	version     = "4.1.0"
	wipSuffix   = "-dev"
)

func Description() string {
	return description
}

func GitSHA() string {
	return gitSHA
}

func Name() string {
	return name
}

func Source() string {
	return source
}

func Version() string {
	return fmt.Sprintf("%s%s", version, wipSuffix)
}
`
)

func TestThatProjectCurrentVersionIsCorrectForWIPVersions(t *testing.T) {
	expectedVersion := "4.1.0"
	currentVersion, err := getVersionInFile([]byte(FileContentsWIP))
	if err != nil {
		t.Fatal(err)
	}
	if currentVersion != expectedVersion {
		t.Fatalf("Got the wrong version, got %#q, expected %#q", currentVersion, expectedVersion)
	}
}

func TestThatProjectCurrentVersionIsCorrectForReleasedVersions(t *testing.T) {
	expectedVersion := "4.1.0"
	currentVersion, err := getVersionInFile([]byte(FileContentsRelease))
	if err != nil {
		t.Fatal(err)
	}
	if currentVersion != expectedVersion {
		t.Fatalf("Got the wrong version, got %#q, expected %#q", currentVersion, expectedVersion)
	}
}
