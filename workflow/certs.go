package workflow

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

func getCertsFromEnv(fs afero.Fs, workingDirectory, envVarPrefix string) (string, string, string, error) {
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
			return "", "", "", fmt.Errorf("could not find cert var: %v", envVarName)
		}

		certFileData, err := base64.StdEncoding.DecodeString(certData)
		if err != nil {
			return "", "", "", fmt.Errorf("could not decode cert: %v", err)
		}

		filePath := filepath.Join(
			workingDirectory,
			strings.ToLower(envVarPrefix)+certDetail.fileNameSuffix,
		)
		if err := afero.WriteFile(fs, filePath, certFileData, 0644); err != nil {
			return "", "", "", fmt.Errorf("could not write cert: %v", err)
		}

		filePaths = append(filePaths, filePath)
	}

	if len(filePaths) != 3 {
		return "", "", "", fmt.Errorf("incorrect number of certs found")
	}

	return filePaths[0], filePaths[1], filePaths[2], nil
}
