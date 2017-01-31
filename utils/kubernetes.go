package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func TemplateKubernetesResources(kubernetesResourcesDirectoryPath, templatedResourcesDirectoryPath, sha string) error {
	if _, err := os.Stat(templatedResourcesDirectoryPath); os.IsNotExist(err) {
		if err := os.Mkdir(templatedResourcesDirectoryPath, 0755); err != nil {
			return err
		}
	}

	files, err := ioutil.ReadDir(kubernetesResourcesDirectoryPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		contents, err := ioutil.ReadFile(filepath.Join(kubernetesResourcesDirectoryPath, file.Name()))
		if err != nil {
			return err
		}

		templatedContents := strings.Replace(string(contents), "%%DOCKER_TAG%%", sha, -1)

		if err := ioutil.WriteFile(filepath.Join(templatedResourcesDirectoryPath, file.Name()), []byte(templatedContents), 0755); err != nil {
			return err
		}
	}

	return nil
}
