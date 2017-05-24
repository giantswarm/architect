package workflow

import (
	"github.com/juju/errgo"
)

var noCertEnvVarError = errgo.New("could not find cert var")

// IsNoCertEnvVar asserts multipleHelmChartsError.
func IsNoCertEnvVar(err error) bool {
	return errgo.Cause(err) == noCertEnvVarError
}

var decodeCertError = errgo.New("could not decode cert")

// IsDecodeCert asserts decodeCertError
func IsDecodeCert(err error) bool {
	return errgo.Cause(err) == decodeCertError
}

var writeCertError = errgo.New("could not write cert")

// IsWriteCert asserts writeCertError
func IsWriteCert(err error) bool {
	return errgo.Cause(err) == writeCertError
}

var incorrectNumberCertsError = errgo.New("incorrect number of certs found")

// IsIncorrectNumberCerts asserts incorrectNumberCertsError
func IsIncorrectNumberCerts(err error) bool {
	return errgo.Cause(err) == incorrectNumberCertsError
}

var emptyWorkingDirectoryError = errgo.New("working directory cannot be empty")

// IsEmptyWorkingDirectory asserts emptyWorkingDirectoryError
func IsEmptyWorkingDirectory(err error) bool {
	return errgo.Cause(err) == emptyWorkingDirectoryError
}

var emptyOrganisationError = errgo.New("organisation cannot be empty")

// IsEmptyOrganisation asserts emptyOrganisationError
func IsEmptyOrganisation(err error) bool {
	return errgo.Cause(err) == emptyOrganisationError
}

var emptyProjectError = errgo.New("project cannot be empty")

// IsEmptyProject asserts emptyProjectError
func IsEmptyProject(err error) bool {
	return errgo.Cause(err) == emptyProjectError
}

var emptyShaError = errgo.New("sha cannot be empty")

// IsEmptySha asserts emptyShaError
func IsEmptySha(err error) bool {
	return errgo.Cause(err) == emptyShaError
}

var emptyRegistryError = errgo.New("registry cannot be empty")

// IsEmptyRegistry asserts emptyRegistryError
func IsEmptyRegistry(err error) bool {
	return errgo.Cause(err) == emptyRegistryError
}

var emptyDockerUsernameError = errgo.New("docker username cannot be empty")

// IsEmptyDockerUsername asserts emptyDockerUsernameError
func IsEmptyDockerUsername(err error) bool {
	return errgo.Cause(err) == emptyDockerUsernameError
}

var emptyDockerPasswordError = errgo.New("docker password cannot be empty")

// IsEmptyDockerPassword asserts emptyDockerPasswordError
func IsEmptyDockerPassword(err error) bool {
	return errgo.Cause(err) == emptyDockerPasswordError
}

var emptyGoosError = errgo.New("goos cannot be empty")

// IsEmptyGoos asserts emptyGoosError
func IsEmptyGoos(err error) bool {
	return errgo.Cause(err) == emptyGoosError
}

var emptyGoarchError = errgo.New("goarch cannot be empty")

// IsEmptyGoarch asserts emptyGoarchError
func IsEmptyGoarch(err error) bool {
	return errgo.Cause(err) == emptyGoarchError
}

var emptyGolangImageError = errgo.New("golang image cannot be empty")

// IsEmptyGolangImage asserts emptyGolangImageError
func IsEmptyGolangImage(err error) bool {
	return errgo.Cause(err) == emptyGolangImageError
}

var emptyGolangVersionError = errgo.New("golang version cannot be empty")

// IsEmptyGolangVersion asserts emptyGolangVersionError
func IsEmptyGolangVersion(err error) bool {
	return errgo.Cause(err) == emptyGolangVersionError
}

var noHelmDirectoryError = errgo.New("cannot find helm directory")

// IsNoHelmDirectory asserts noHelmDirectoryError
func IsNoHelmDirectory(err error) bool {
	return errgo.Cause(err) == noHelmDirectoryError
}

var emptyKubernetesAPIServerError = errgo.New("kubernetes api server cannot be empty")

// IsEmptyKubernetesAPIServer asserts emptyKubernetesAPIServerError
func IsEmptyKubernetesAPIServer(err error) bool {
	return errgo.Cause(err) == emptyKubernetesAPIServerError
}

var emptyKubernetesCaPathError = errgo.New("kubernetes ca path cannot be empty")

// IsEmptyKubernetesCaPath asserts emptyKubernetesCaPathError
func IsEmptyKubernetesCaPath(err error) bool {
	return errgo.Cause(err) == emptyKubernetesCaPathError
}

var emptyKubernetesCrtPathError = errgo.New("kubernetes crt path cannot be empty")

// IsEmptyKubernetesCrtPath asserts emptyKubernetesCAPathError
func IsEmptyKubernetesCrtPath(err error) bool {
	return errgo.Cause(err) == emptyKubernetesCrtPathError
}

var emptyKubernetesKeyPathError = errgo.New("kubernetes key path cannot be empty")

// IsEmptyKubernetesKeyPath asserts emptyKubernetesKeyPathError
func IsEmptyKubernetesCAPath(err error) bool {
	return errgo.Cause(err) == emptyKubernetesKeyPathError
}

var emptyKubectlVersionError = errgo.New("kubectl version cannot be empty")

// IsEmptyKubectlVersion asserts emptyKubectlVersionError
func IsEmptyKubectlVersion(err error) bool {
	return errgo.Cause(err) == emptyKubectlVersionError
}

var emptyKubernetesResourcesDirectoryPath = errgo.New("kubenrnetes resources directory path cannot be empty")

// IsEmptyKubernetesResourcesDirectoryPath asserts emptyKubernetesResourcesDirectoryPath
func IsEmptyKubernetesResourcesDirectoryPath(err error) bool {
	return errgo.Cause(err) == emptyKubernetesResourcesDirectoryPath
}
