package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestTemplateKubernetesResources(t *testing.T) {
	var filePermission os.FileMode = 0644
	var directoryPermission os.FileMode = 0644

	resourcesPath := "./kubernetes/"
	templatedResourcesPath := "./kubernetes-templated/"

	sha := "1cd72a25e16e93da14f08d95bd98662f8827028e" // random, no specific meaning
	testData := []byte("this is some test data")

	tests := []struct {
		setUp func(afero.Fs) error
		check func(afero.Fs) error
	}{
		// Test an empty resources directory produces an empty templates resources directory
		{
			setUp: func(fs afero.Fs) error {
				if err := fs.Mkdir(resourcesPath, directoryPermission); err != nil {
					return err
				}

				return nil
			},
			check: func(fs afero.Fs) error {
				fileInfos, err := afero.ReadDir(fs, templatedResourcesPath)
				if err != nil {
					return err
				}

				if len(fileInfos) != 0 {
					return fmt.Errorf("multiple files found in templated resources directory")
				}

				return nil
			},
		},

		// Test a resources directory with a deployment is templated correctly
		{
			setUp: func(fs afero.Fs) error {
				if err := fs.Mkdir(resourcesPath, directoryPermission); err != nil {
					return err
				}

				if err := afero.WriteFile(
					fs,
					filepath.Join(resourcesPath, "deployment.yml"),
					[]byte(`%%DOCKER_TAG%%`),
					filePermission,
				); err != nil {
					return err
				}

				return nil
			},
			check: func(fs afero.Fs) error {
				fileInfos, err := afero.ReadDir(fs, templatedResourcesPath)
				if err != nil {
					return err
				}

				if len(fileInfos) != 1 {
					return fmt.Errorf("did not find only one file in templated resources directory")
				}

				if fileInfos[0].Name() != "deployment.yml" {
					return fmt.Errorf("deployment not found in templates resources directory")
				}

				bytes, err := afero.ReadFile(fs, filepath.Join(templatedResourcesPath, "deployment.yml"))
				if err != nil {
					return err
				}
				if !strings.Contains(string(bytes), sha) {
					return fmt.Errorf("sha not found in deployment")
				}

				return nil
			},
		},

		// Test a nested resources directory is handled correctly
		{
			setUp: func(fs afero.Fs) error {
				if err := fs.Mkdir(resourcesPath, directoryPermission); err != nil {
					return err
				}

				if err := fs.Mkdir(filepath.Join(resourcesPath, "foo/"), directoryPermission); err != nil {
					return err
				}

				if err := afero.WriteFile(
					fs,
					filepath.Join(resourcesPath, "foo/", "deployment.yml"),
					testData,
					filePermission,
				); err != nil {
					return err
				}

				if err := fs.Mkdir(filepath.Join(resourcesPath, "bar/"), directoryPermission); err != nil {
					return err
				}

				if err := afero.WriteFile(
					fs,
					filepath.Join(resourcesPath, "bar/", "deployment.yml"),
					testData,
					filePermission,
				); err != nil {
					return err
				}

				return nil
			},
			check: func(fs afero.Fs) error {
				fileInfos, err := afero.ReadDir(fs, templatedResourcesPath)
				if err != nil {
					return err
				}

				if len(fileInfos) != 2 || !fileInfos[0].IsDir() || !fileInfos[1].IsDir() {
					return fmt.Errorf("did not find two directories in template resources directory")
				}

				fooContents, err := afero.ReadFile(
					fs,
					filepath.Join(resourcesPath, "foo", "deployment.yml"),
				)
				if string(fooContents) != string(testData) {
					return fmt.Errorf("foo file data did not match test data")
				}

				barContents, err := afero.ReadFile(
					fs,
					filepath.Join(resourcesPath, "bar", "deployment.yml"),
				)
				if string(barContents) != string(testData) {
					return fmt.Errorf("bar file data did not match test data")
				}

				return nil
			},
		},
	}

	for index, test := range tests {
		fs := afero.NewMemMapFs()

		if err := test.setUp(fs); err != nil {
			t.Fatalf("%v: unexpected error during setup: %v\n", index, err)
		}

		if err := TemplateKubernetesResources(fs, resourcesPath, templatedResourcesPath, sha); err != nil {
			t.Fatalf("%v: unexpected error during templating: %v\n", index, err)
		}

		if err := test.check(fs); err != nil {
			t.Fatalf("%v: unexpected error during check: %v\n", index, err)
		}
	}
}
