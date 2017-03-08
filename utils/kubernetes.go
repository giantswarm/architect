package utils

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

func TemplateFile(fs afero.Fs, path, templatedPath, sha string) error {
	contents, err := afero.ReadFile(fs, path)
	if err != nil {
		return err
	}

	templatedContents := strings.Replace(string(contents), "%%DOCKER_TAG%%", sha, -1)

	if err := afero.WriteFile(fs, templatedPath, []byte(templatedContents), 0644); err != nil {
		return err
	}

	return nil
}

func TemplateKubernetesResources(fs afero.Fs, resourcesDir, templatesDir, sha string) error {
	// Generate a list of files to template
	type FileToTemplate struct {
		path               string
		templatedDirectory string
		templatedPath      string
	}

	filesToTemplate := []FileToTemplate{}

	fileInfos, err := afero.ReadDir(fs, resourcesDir)
	if err != nil {
		return err
	}

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			subFileInfos, err := afero.ReadDir(fs, filepath.Join(resourcesDir, fileInfo.Name()))
			if err != nil {
				return err
			}

			for _, subFileInfo := range subFileInfos {
				if !subFileInfo.IsDir() {
					filesToTemplate = append(filesToTemplate, FileToTemplate{
						path:               filepath.Join(resourcesDir, fileInfo.Name(), subFileInfo.Name()),
						templatedDirectory: filepath.Join(templatesDir, fileInfo.Name()),
						templatedPath:      filepath.Join(templatesDir, fileInfo.Name(), subFileInfo.Name()),
					})
				}
			}
		} else {
			filesToTemplate = append(filesToTemplate, FileToTemplate{
				path:          filepath.Join(resourcesDir, fileInfo.Name()),
				templatedPath: filepath.Join(templatesDir, fileInfo.Name()),
			})
		}
	}

	// Create the templated resources directory, if it does not exist
	if _, err := fs.Stat(templatesDir); os.IsNotExist(err) {
		if err := fs.Mkdir(templatesDir, 0644); err != nil {
			return err
		}
	}

	// And template the files
	for _, fileToTemplate := range filesToTemplate {
		// Create the directory, if it does not exist
		if fileToTemplate.templatedDirectory != "" {
			if _, err := fs.Stat(fileToTemplate.templatedDirectory); os.IsNotExist(err) {
				if err := fs.Mkdir(fileToTemplate.templatedDirectory, 0644); err != nil {
					return err
				}
			}
		}

		if err := TemplateFile(fs, fileToTemplate.path, fileToTemplate.templatedPath, sha); err != nil {
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
