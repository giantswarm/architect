package utils

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

func TemplateKubernetesResources(fs afero.Fs, kubernetesResourcesDirectoryPath, templatedResourcesDirectoryPath, sha string) error {
	if _, err := fs.Stat(templatedResourcesDirectoryPath); os.IsNotExist(err) {
		if err := fs.Mkdir(templatedResourcesDirectoryPath, 0755); err != nil {
			return err
		}
	}

	files, err := afero.ReadDir(fs, kubernetesResourcesDirectoryPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		contents, err := afero.ReadFile(fs, filepath.Join(kubernetesResourcesDirectoryPath, file.Name()))
		if err != nil {
			return err
		}

		templatedContents := strings.Replace(string(contents), "%%DOCKER_TAG%%", sha, -1)

		if err := afero.WriteFile(fs, filepath.Join(templatedResourcesDirectoryPath, file.Name()), []byte(templatedContents), 0755); err != nil {
			return err
		}
	}

	return nil
}

// K8SCertsFromEnv attempts to load certificates from the environment
func K8SCertsFromEnv(fs afero.Fs, workingDirectory string) (string, string, string, error) {
	caName := "ca.pem"
	crtName := "crt.pem"
	keyName := "key.pem"

	assets := []struct {
		envVar string
		name   string
	}{
		{envVar: "G8S_CA", name: caName},
		{envVar: "G8S_CRT", name: crtName},
		{envVar: "G8S_KEY", name: keyName},
	}

	for _, asset := range assets {
		encodedData := os.Getenv(asset.envVar)
		if encodedData == "" {
			return "", "", "", fmt.Errorf("could not load certificate data from %v\n", asset.envVar)
		}

		data, err := base64.StdEncoding.DecodeString(encodedData)
		if err != nil {
			return "", "", "", err
		}

		if err := afero.WriteFile(fs, asset.name, data, 0644); err != nil {
			return "", "", "", err
		}
	}

	caPath := filepath.Join(workingDirectory, caName)
	crtPath := filepath.Join(workingDirectory, crtName)
	keyPath := filepath.Join(workingDirectory, keyName)

	return caPath, crtPath, keyPath, nil
}
