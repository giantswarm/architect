package utils

import (
	"os"
	"strings"

	"github.com/spf13/afero"
)

func TemplateFile(fs afero.Fs, path, sha, ingressTag string) error {
	contents, err := afero.ReadFile(fs, path)
	if err != nil {
		return err
	}

	templatedContents := strings.Replace(string(contents), "%%DOCKER_TAG%%", sha, -1)
	templatedContents = strings.Replace(string(templatedContents), "%%INGRESS_TAG%%", ingressTag, -1)

	if err := afero.WriteFile(fs, path, []byte(templatedContents), 0644); err != nil {
		return err
	}

	return nil
}

func TemplateKubernetesResources(fs afero.Fs, resourcesDir, sha, ingressTag string) error {
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if err := TemplateFile(fs, path, sha, ingressTag); err != nil {
			return err
		}

		return nil
	}

	if err := afero.Walk(fs, resourcesDir, walkFunc); err != nil {
		return err
	}

	return nil
}
