package workflow

import (
	"github.com/giantswarm/microerror"
)

var noCertEnvVarError = &microerror.Error{
	Kind: "noCertEnvVarError",
}

// IsNoCertEnvVar asserts noCertEnvVarError
func IsNoCertEnvVar(err error) bool {
	return microerror.Cause(err) == noCertEnvVarError
}

var decodeCertError = &microerror.Error{
	Kind: "decodeCertError",
}

// IsDecodeCert asserts decodeCertError
func IsDecodeCert(err error) bool {
	return microerror.Cause(err) == decodeCertError
}

var writeCertError = &microerror.Error{
	Kind: "writeCertError",
}

// IsWriteCert asserts writeCertError
func IsWriteCert(err error) bool {
	return microerror.Cause(err) == writeCertError
}

var incorrectNumberCertsError = &microerror.Error{
	Kind: "incorrectNumberCertsError",
}

// IsIncorrectNumberCerts asserts incorrectNumberCertsError
func IsIncorrectNumberCerts(err error) bool {
	return microerror.Cause(err) == incorrectNumberCertsError
}

var emptyChannelError = &microerror.Error{
	Kind: "emptyChannelError",
}

// IsEmptyChannel asserts emptyChannelError
func IsEmptyChannel(err error) bool {
	return microerror.Cause(err) == emptyChannelError
}

var emptyWorkingDirectoryError = &microerror.Error{
	Kind: "emptyWorkingDirectoryError",
}

// IsEmptyWorkingDirectory asserts emptyWorkingDirectoryError
func IsEmptyWorkingDirectory(err error) bool {
	return microerror.Cause(err) == emptyWorkingDirectoryError
}

var emptyOrganisationError = &microerror.Error{
	Kind: "emptyOrganisationError",
}

// IsEmptyOrganisation asserts emptyOrganisationError
func IsEmptyOrganisation(err error) bool {
	return microerror.Cause(err) == emptyOrganisationError
}

var emptyProjectError = &microerror.Error{
	Kind: "emptyProjectError",
}

// IsEmptyProject asserts emptyProjectError
func IsEmptyProject(err error) bool {
	return microerror.Cause(err) == emptyProjectError
}

var emptyRefError = &microerror.Error{
	Kind: "emptyRefError",
}

// IsEmptyRef asserts emptyRefError
func IsEmptyRef(err error) bool {
	return microerror.Cause(err) == emptyRefError
}

var emptyShaError = &microerror.Error{
	Kind: "emptyShaError",
}

// IsEmptySha asserts emptyShaError
func IsEmptySha(err error) bool {
	return microerror.Cause(err) == emptyShaError
}

var emptyRegistryError = &microerror.Error{
	Kind: "emptyRegistryError",
}

// IsEmptyRegistry asserts emptyRegistryError
func IsEmptyRegistry(err error) bool {
	return microerror.Cause(err) == emptyRegistryError
}

var emptyDockerUsernameError = &microerror.Error{
	Kind: "emptyDockerUsernameError",
}

// IsEmptyDockerUsername asserts emptyDockerUsernameError
func IsEmptyDockerUsername(err error) bool {
	return microerror.Cause(err) == emptyDockerUsernameError
}

var emptyDockerPasswordError = &microerror.Error{
	Kind: "emptyDockerPasswordError",
}

// IsEmptyDockerPassword asserts emptyDockerPasswordError
func IsEmptyDockerPassword(err error) bool {
	return microerror.Cause(err) == emptyDockerPasswordError
}

var emptyGoosError = &microerror.Error{
	Kind: "emptyGoosError",
}

// IsEmptyGoos asserts emptyGoosError
func IsEmptyGoos(err error) bool {
	return microerror.Cause(err) == emptyGoosError
}

var emptyGoarchError = &microerror.Error{
	Kind: "emptyGoarchError",
}

// IsEmptyGoarch asserts emptyGoarchError
func IsEmptyGoarch(err error) bool {
	return microerror.Cause(err) == emptyGoarchError
}

var emptyGolangImageError = &microerror.Error{
	Kind: "emptyGolangImageError",
}

// IsEmptyGolangImage asserts emptyGolangImageError
func IsEmptyGolangImage(err error) bool {
	return microerror.Cause(err) == emptyGolangImageError
}

var emptyGolangVersionError = &microerror.Error{
	Kind: "emptyGolangVersionError",
}

// IsEmptyGolangVersion asserts emptyGolangVersionError
func IsEmptyGolangVersion(err error) bool {
	return microerror.Cause(err) == emptyGolangVersionError
}

var failedExecutionError = &microerror.Error{
	Kind: "failedExecutionError",
}

// IsFailedExecution asserts failedExecutionError.
func IsFailedExecutionError(err error) bool {
	return microerror.Cause(err) == failedExecutionError
}

var noHelmDirectoryError = &microerror.Error{
	Kind: "noHelmDirectoryError",
}

// IsNoHelmDirectory asserts noHelmDirectoryError
func IsNoHelmDirectory(err error) bool {
	return microerror.Cause(err) == noHelmDirectoryError
}

var emptyKubernetesAPIServerError = &microerror.Error{
	Kind: "emptyKubernetesAPIServerError",
}

// IsEmptyKubernetesAPIServer asserts emptyKubernetesAPIServerError
func IsEmptyKubernetesAPIServer(err error) bool {
	return microerror.Cause(err) == emptyKubernetesAPIServerError
}

var emptyKubernetesCaPathError = &microerror.Error{
	Kind: "emptyKubernetesCaPathError",
}

// IsEmptyKubernetesCaPath asserts emptyKubernetesCaPathError
func IsEmptyKubernetesCaPath(err error) bool {
	return microerror.Cause(err) == emptyKubernetesCaPathError
}

var emptyKubernetesCrtPathError = &microerror.Error{
	Kind: "emptyKubernetesCrtPathError",
}

// IsEmptyKubernetesCrtPath asserts emptyKubernetesCAPathError
func IsEmptyKubernetesCrtPath(err error) bool {
	return microerror.Cause(err) == emptyKubernetesCrtPathError
}

var emptyKubernetesKeyPathError = &microerror.Error{
	Kind: "emptyKubernetesKeyPathError",
}

// IsEmptyKubernetesKeyPath asserts emptyKubernetesKeyPathError
func IsEmptyKubernetesCAPath(err error) bool {
	return microerror.Cause(err) == emptyKubernetesKeyPathError
}

var emptyKubectlVersionError = &microerror.Error{
	Kind: "emptyKubectlVersionError",
}

// IsEmptyKubectlVersion asserts emptyKubectlVersionError
func IsEmptyKubectlVersion(err error) bool {
	return microerror.Cause(err) == emptyKubectlVersionError
}

var invalidHelmDirectoryError = &microerror.Error{
	Kind: "invalidHelmDirectoryError",
}

// IsInvalidHelmDirectory asserts invalidHelmDirectoryError.
func IsInvalidHelmDirectory(err error) bool {
	return microerror.Cause(err) == invalidHelmDirectoryError
}

var missingFileError = &microerror.Error{
	Kind: "missingFileError",
}

// IsMissingFileError asserts missingFileError.
func IsMissingFileError(err error) bool {
	return microerror.Cause(err) == missingFileError
}
