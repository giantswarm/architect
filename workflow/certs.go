package workflow

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"

	"github.com/giantswarm/microerror"
)

func CertsFromEnv(fs afero.Fs, workingDirectory, envVarPrefix string) (string, string, string, error) {
	certDetails := []struct {
		envVarSuffix   string
		fileNameSuffix string
	}{
		{envVarSuffix: "_CA", fileNameSuffix: "-ca.pem"},
		{envVarSuffix: "_CRT", fileNameSuffix: "-crt.pem"},
		{envVarSuffix: "_KEY", fileNameSuffix: "-key.pem"},
	}

	filePaths := []string{}

	for _, certDetail := range certDetails {
		envVarName := envVarPrefix + certDetail.envVarSuffix

		certData := os.Getenv(envVarName)
		if certData == "" {
			return "", "", "", microerror.Maskf(noCertEnvVarError, envVarName)
		}

		certFileData, err := base64.StdEncoding.DecodeString(certData)
		if err != nil {
			return "", "", "", microerror.Maskf(decodeCertError, err.Error())
		}

		filePath := filepath.Join(
			workingDirectory,
			strings.ToLower(envVarPrefix)+certDetail.fileNameSuffix,
		)
		if err := afero.WriteFile(fs, filePath, certFileData, 0644); err != nil {
			return "", "", "", microerror.Maskf(writeCertError, err.Error())
		}

		filePaths = append(filePaths, filePath)
	}

	if len(filePaths) != 3 {
		return "", "", "", microerror.Maskf(incorrectNumberCertsError, string(len(filePaths)))
	}

	return filePaths[0], filePaths[1], filePaths[2], nil
}
