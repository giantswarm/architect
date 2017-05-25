package workflow

import (
	"github.com/juju/errgo"
)

var noCertEnvVarError = errgo.New("no cert env var")

// IsNoCertEnvVar asserts multipleHelmChartsError.
func IsNoCertEnvVar(err error) bool {
	return errgo.Cause(err) == noCertEnvVarError
}

var decodeCertError = errgo.New("decode cert")

// IsDecodeCert asserts decodeCertError
func IsDecodeCert(err error) bool {
	return errgo.Cause(err) == decodeCertError
}

var writeCertError = errgo.New("write cert")

// IsWriteCert asserts writeCertError
func IsWriteCert(err error) bool {
	return errgo.Cause(err) == writeCertError
}

var incorrectNumberCertsError = errgo.New("incorrect number certs")

// IsIncorrectNumberCerts asserts incorrectNumberCertsError
func IsIncorrectNumberCerts(err error) bool {
	return errgo.Cause(err) == incorrectNumberCertsError
}

var emptyWorkingDirectoryError = errgo.New("empty working directory")

// IsEmptyWorkingDirectory asserts emptyWorkingDirectoryError
func IsEmptyWorkingDirectory(err error) bool {
	return errgo.Cause(err) == emptyWorkingDirectoryError
}

var emptyOrganisationError = errgo.New("empty organisation")

// IsEmptyOrganisation asserts emptyOrganisationError
func IsEmptyOrganisation(err error) bool {
	return errgo.Cause(err) == emptyOrganisationError
}

var emptyProjectError = errgo.New("empty project")

// IsEmptyProject asserts emptyProjectError
func IsEmptyProject(err error) bool {
	return errgo.Cause(err) == emptyProjectError
}

var emptyShaError = errgo.New("empty sha")

// IsEmptySha asserts emptyShaError
func IsEmptySha(err error) bool {
	return errgo.Cause(err) == emptyShaError
}

var emptyRegistryError = errgo.New("empty registry")

// IsEmptyRegistry asserts emptyRegistryError
func IsEmptyRegistry(err error) bool {
	return errgo.Cause(err) == emptyRegistryError
}

var emptyDockerUsernameError = errgo.New("empty docker username")

// IsEmptyDockerUsername asserts emptyDockerUsernameError
func IsEmptyDockerUsername(err error) bool {
	return errgo.Cause(err) == emptyDockerUsernameError
}

var emptyDockerPasswordError = errgo.New("empty docker password")

// IsEmptyDockerPassword asserts emptyDockerPasswordError
func IsEmptyDockerPassword(err error) bool {
	return errgo.Cause(err) == emptyDockerPasswordError
}

var emptyGoosError = errgo.New("empty goos")

// IsEmptyGoos asserts emptyGoosError
func IsEmptyGoos(err error) bool {
	return errgo.Cause(err) == emptyGoosError
}

var emptyGoarchError = errgo.New("empty goarch")

// IsEmptyGoarch asserts emptyGoarchError
func IsEmptyGoarch(err error) bool {
	return errgo.Cause(err) == emptyGoarchError
}

var emptyGolangImageError = errgo.New("empty golang image")

// IsEmptyGolangImage asserts emptyGolangImageError
func IsEmptyGolangImage(err error) bool {
	return errgo.Cause(err) == emptyGolangImageError
}

var emptyGolangVersionError = errgo.New("empty golang version")

// IsEmptyGolangVersion asserts emptyGolangVersionError
func IsEmptyGolangVersion(err error) bool {
	return errgo.Cause(err) == emptyGolangVersionError
}

var noHelmDirectoryError = errgo.New("no helm directory")

// IsNoHelmDirectory asserts noHelmDirectoryError
func IsNoHelmDirectory(err error) bool {
	return errgo.Cause(err) == noHelmDirectoryError
}

var emptyKubernetesAPIServerError = errgo.New("empty kubernetes api server")

// IsEmptyKubernetesAPIServer asserts emptyKubernetesAPIServerError
func IsEmptyKubernetesAPIServer(err error) bool {
	return errgo.Cause(err) == emptyKubernetesAPIServerError
}

var emptyKubernetesCaPathError = errgo.New("empty kubernetes ca path")

// IsEmptyKubernetesCaPath asserts emptyKubernetesCaPathError
func IsEmptyKubernetesCaPath(err error) bool {
	return errgo.Cause(err) == emptyKubernetesCaPathError
}

var emptyKubernetesCrtPathError = errgo.New("empty kubernetes crt path")

// IsEmptyKubernetesCrtPath asserts emptyKubernetesCAPathError
func IsEmptyKubernetesCrtPath(err error) bool {
	return errgo.Cause(err) == emptyKubernetesCrtPathError
}

var emptyKubernetesKeyPathError = errgo.New("empty kubernetes key path")

// IsEmptyKubernetesKeyPath asserts emptyKubernetesKeyPathError
func IsEmptyKubernetesCAPath(err error) bool {
	return errgo.Cause(err) == emptyKubernetesKeyPathError
}

var emptyKubectlVersionError = errgo.New("empty kubectl version")

// IsEmptyKubectlVersion asserts emptyKubectlVersionError
func IsEmptyKubectlVersion(err error) bool {
	return errgo.Cause(err) == emptyKubectlVersionError
}

var emptyKubernetesResourcesDirectoryPath = errgo.New("empty kubernetes resources directory path")

// IsEmptyKubernetesResourcesDirectoryPath asserts emptyKubernetesResourcesDirectoryPath
func IsEmptyKubernetesResourcesDirectoryPath(err error) bool {
	return errgo.Cause(err) == emptyKubernetesResourcesDirectoryPath
}
