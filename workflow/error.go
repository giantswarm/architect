package workflow

import (
	"github.com/giantswarm/microerror"
)

var noCertEnvVarError = microerror.New("no cert env var")

// IsNoCertEnvVar asserts noCertEnvVarError
func IsNoCertEnvVar(err error) bool {
	return microerror.Cause(err) == noCertEnvVarError
}

var noGolangPackagesError = microerror.New("no golang packages found")

// IsNoGolangPackages asserts noGolangPackagesError
func IsNoGolangPackages(err error) bool {
	return microerror.Cause(err) == noGolangPackagesError
}

var decodeCertError = microerror.New("decode cert")

// IsDecodeCert asserts decodeCertError
func IsDecodeCert(err error) bool {
	return microerror.Cause(err) == decodeCertError
}

var writeCertError = microerror.New("write cert")

// IsWriteCert asserts writeCertError
func IsWriteCert(err error) bool {
	return microerror.Cause(err) == writeCertError
}

var incorrectNumberCertsError = microerror.New("incorrect number certs")

// IsIncorrectNumberCerts asserts incorrectNumberCertsError
func IsIncorrectNumberCerts(err error) bool {
	return microerror.Cause(err) == incorrectNumberCertsError
}

var emptyChannelError = microerror.New("empty channel")

// IsEmptyChannel asserts emptyChannelError
func IsEmptyChannel(err error) bool {
	return microerror.Cause(err) == emptyChannelError
}

var emptyWorkingDirectoryError = microerror.New("empty working directory")

// IsEmptyWorkingDirectory asserts emptyWorkingDirectoryError
func IsEmptyWorkingDirectory(err error) bool {
	return microerror.Cause(err) == emptyWorkingDirectoryError
}

var emptyOrganisationError = microerror.New("empty organisation")

// IsEmptyOrganisation asserts emptyOrganisationError
func IsEmptyOrganisation(err error) bool {
	return microerror.Cause(err) == emptyOrganisationError
}

var emptyProjectError = microerror.New("empty project")

// IsEmptyProject asserts emptyProjectError
func IsEmptyProject(err error) bool {
	return microerror.Cause(err) == emptyProjectError
}

var emptyShaError = microerror.New("empty sha")

// IsEmptySha asserts emptyShaError
func IsEmptySha(err error) bool {
	return microerror.Cause(err) == emptyShaError
}

var emptyRegistryError = microerror.New("empty registry")

// IsEmptyRegistry asserts emptyRegistryError
func IsEmptyRegistry(err error) bool {
	return microerror.Cause(err) == emptyRegistryError
}

var emptyDockerUsernameError = microerror.New("empty docker username")

// IsEmptyDockerUsername asserts emptyDockerUsernameError
func IsEmptyDockerUsername(err error) bool {
	return microerror.Cause(err) == emptyDockerUsernameError
}

var emptyDockerPasswordError = microerror.New("empty docker password")

// IsEmptyDockerPassword asserts emptyDockerPasswordError
func IsEmptyDockerPassword(err error) bool {
	return microerror.Cause(err) == emptyDockerPasswordError
}

var emptyGoosError = microerror.New("empty goos")

// IsEmptyGoos asserts emptyGoosError
func IsEmptyGoos(err error) bool {
	return microerror.Cause(err) == emptyGoosError
}

var emptyGoarchError = microerror.New("empty goarch")

// IsEmptyGoarch asserts emptyGoarchError
func IsEmptyGoarch(err error) bool {
	return microerror.Cause(err) == emptyGoarchError
}

var emptyGolangImageError = microerror.New("empty golang image")

// IsEmptyGolangImage asserts emptyGolangImageError
func IsEmptyGolangImage(err error) bool {
	return microerror.Cause(err) == emptyGolangImageError
}

var emptyGolangVersionError = microerror.New("empty golang version")

// IsEmptyGolangVersion asserts emptyGolangVersionError
func IsEmptyGolangVersion(err error) bool {
	return microerror.Cause(err) == emptyGolangVersionError
}

var noHelmDirectoryError = microerror.New("no helm directory")

// IsNoHelmDirectory asserts noHelmDirectoryError
func IsNoHelmDirectory(err error) bool {
	return microerror.Cause(err) == noHelmDirectoryError
}

var emptyKubernetesAPIServerError = microerror.New("empty kubernetes api server")

// IsEmptyKubernetesAPIServer asserts emptyKubernetesAPIServerError
func IsEmptyKubernetesAPIServer(err error) bool {
	return microerror.Cause(err) == emptyKubernetesAPIServerError
}

var emptyKubernetesCaPathError = microerror.New("empty kubernetes ca path")

// IsEmptyKubernetesCaPath asserts emptyKubernetesCaPathError
func IsEmptyKubernetesCaPath(err error) bool {
	return microerror.Cause(err) == emptyKubernetesCaPathError
}

var emptyKubernetesCrtPathError = microerror.New("empty kubernetes crt path")

// IsEmptyKubernetesCrtPath asserts emptyKubernetesCAPathError
func IsEmptyKubernetesCrtPath(err error) bool {
	return microerror.Cause(err) == emptyKubernetesCrtPathError
}

var emptyKubernetesKeyPathError = microerror.New("empty kubernetes key path")

// IsEmptyKubernetesKeyPath asserts emptyKubernetesKeyPathError
func IsEmptyKubernetesCAPath(err error) bool {
	return microerror.Cause(err) == emptyKubernetesKeyPathError
}

var emptyKubectlVersionError = microerror.New("empty kubectl version")

// IsEmptyKubectlVersion asserts emptyKubectlVersionError
func IsEmptyKubectlVersion(err error) bool {
	return microerror.Cause(err) == emptyKubectlVersionError
}

var invalidHelmDirectoryError = microerror.New("invalid helm directory")

// IsInvalidHelmDirectory asserts invalidHelmDirectoryError.
func IsInvalidHelmDirectory(err error) bool {
	return microerror.Cause(err) == invalidHelmDirectoryError
}
