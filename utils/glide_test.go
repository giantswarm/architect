package utils

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spf13/afero"

	"github.com/giantswarm/microerror"
)

func TestNoVendor(t *testing.T) {
	var directoryPermission os.FileMode = 0644

	tests := []struct {
		workingDirectory    string
		setUp               func(afero.Fs, string) error
		expectedDirectories []string
	}{
		// Test a completely empty working directory
		{
			workingDirectory: "/",
			setUp: func(fs afero.Fs, wd string) error {
				return nil
			},
			expectedDirectories: []string{},
		},

		// Test one empty directory
		{
			workingDirectory: "/",
			setUp: func(fs afero.Fs, wd string) error {
				if err := fs.Mkdir(filepath.Join(wd, "pkg"), directoryPermission); err != nil {
					return microerror.Mask(err)
				}

				return nil
			},
			expectedDirectories: []string{},
		},

		// Test one golang file
		{
			workingDirectory: "/",
			setUp: func(fs afero.Fs, wd string) error {
				if _, err := fs.Create(filepath.Join(wd, "main.go")); err != nil {
					return microerror.Mask(err)
				}

				return nil
			},
			expectedDirectories: []string{"."},
		},

		// Test one directory with a golang file in
		{
			workingDirectory: "/",
			setUp: func(fs afero.Fs, wd string) error {
				if err := fs.Mkdir(filepath.Join(wd, "pkg"), directoryPermission); err != nil {
					return microerror.Mask(err)
				}

				if _, err := fs.Create(filepath.Join(wd, "pkg", "main.go")); err != nil {
					return microerror.Mask(err)
				}

				return nil
			},
			expectedDirectories: []string{"./pkg/..."},
		},

		// Test two directories with one golang file in each
		{
			workingDirectory: "/",
			setUp: func(fs afero.Fs, wd string) error {
				if err := fs.Mkdir(filepath.Join(wd, "bar"), directoryPermission); err != nil {
					return microerror.Mask(err)
				}

				if _, err := fs.Create(filepath.Join(wd, "bar", "main.go")); err != nil {
					return microerror.Mask(err)
				}

				if err := fs.Mkdir(filepath.Join(wd, "foo"), directoryPermission); err != nil {
					return microerror.Mask(err)
				}

				if _, err := fs.Create(filepath.Join(wd, "foo", "main.go")); err != nil {
					return microerror.Mask(err)
				}

				return nil
			},
			expectedDirectories: []string{"./bar/...", "./foo/..."},
		},

		// Test one directory with two golang files in
		{
			workingDirectory: "/",
			setUp: func(fs afero.Fs, wd string) error {
				if err := fs.Mkdir(filepath.Join(wd, "pkg"), directoryPermission); err != nil {
					return microerror.Mask(err)
				}

				if _, err := fs.Create(filepath.Join(wd, "pkg", "main.go")); err != nil {
					return microerror.Mask(err)
				}

				if _, err := fs.Create(filepath.Join(wd, "pkg", "error.go")); err != nil {
					return microerror.Mask(err)
				}

				return nil
			},
			expectedDirectories: []string{"./pkg/..."},
		},

		// Test one directory with a golang file in, with an alternative working directory
		{
			workingDirectory: "/home/ubuntu/api",
			setUp: func(fs afero.Fs, wd string) error {
				if err := fs.Mkdir(filepath.Join(wd, "pkg"), directoryPermission); err != nil {
					return microerror.Mask(err)
				}

				if _, err := fs.Create(filepath.Join(wd, "pkg", "main.go")); err != nil {
					return microerror.Mask(err)
				}

				return nil
			},
			expectedDirectories: []string{"./pkg/..."},
		},

		// Test one directory, with a sub directory containing a golang file
		{
			workingDirectory: "/",
			setUp: func(fs afero.Fs, wd string) error {
				if err := fs.Mkdir(filepath.Join(wd, "pkg"), directoryPermission); err != nil {
					return microerror.Mask(err)
				}

				if err := fs.Mkdir(filepath.Join(wd, "pkg", "subpkg"), directoryPermission); err != nil {
					return microerror.Mask(err)
				}

				if _, err := fs.Create(filepath.Join(wd, "pkg", "subpkg", "main.go")); err != nil {
					return microerror.Mask(err)
				}

				return nil
			},
			expectedDirectories: []string{"./pkg/..."},
		},

		// Test one directory, with two sub directories containing a golang file each
		{
			workingDirectory: "/",
			setUp: func(fs afero.Fs, wd string) error {
				if err := fs.Mkdir(filepath.Join(wd, "pkg"), directoryPermission); err != nil {
					return microerror.Mask(err)
				}

				if err := fs.Mkdir(filepath.Join(wd, "pkg", "subpkg1"), directoryPermission); err != nil {
					return microerror.Mask(err)
				}

				if _, err := fs.Create(filepath.Join(wd, "pkg", "subpkg1", "main.go")); err != nil {
					return microerror.Mask(err)
				}

				if err := fs.Mkdir(filepath.Join(wd, "pkg", "subpkg2"), directoryPermission); err != nil {
					return microerror.Mask(err)
				}

				if _, err := fs.Create(filepath.Join(wd, "pkg", "subpkg2", "main.go")); err != nil {
					return microerror.Mask(err)
				}

				return nil
			},
			expectedDirectories: []string{"./pkg/..."},
		},

		// Test one directory which is ordered lexicographically after the file in the root dir
		{
			workingDirectory: "/",
			setUp: func(fs afero.Fs, wd string) error {
				if _, err := fs.Create(filepath.Join(wd, "main.go")); err != nil {
					return microerror.Mask(err)
				}

				if err := fs.Mkdir(filepath.Join(wd, "zzz"), directoryPermission); err != nil {
					return microerror.Mask(err)
				}

				if _, err := fs.Create(filepath.Join(wd, "zzz", "sleep.go")); err != nil {
					return microerror.Mask(err)
				}

				return nil
			},
			expectedDirectories: []string{".", "./zzz/..."},
		},

		// Test the vendor directory is removed
		{
			workingDirectory: "/",
			setUp: func(fs afero.Fs, wd string) error {
				if err := fs.Mkdir(filepath.Join(wd, "vendor"), directoryPermission); err != nil {
					return microerror.Mask(err)
				}

				if _, err := fs.Create(filepath.Join(wd, "vendor", "vendor.go")); err != nil {
					return microerror.Mask(err)
				}

				return nil
			},
			expectedDirectories: []string{},
		},

		// Test a directory starting with an underscore is ingored
		{
			workingDirectory: "/",
			setUp: func(fs afero.Fs, wd string) error {
				if err := fs.Mkdir(filepath.Join(wd, "_tests"), directoryPermission); err != nil {
					return microerror.Mask(err)
				}

				if _, err := fs.Create(filepath.Join(wd, "_tests", "test.go")); err != nil {
					return microerror.Mask(err)
				}

				return nil
			},
			expectedDirectories: []string{},
		},

		// Test two directories, containing golang files, a golang file in the root dir,
		// and an alternative working directory
		{
			workingDirectory: "/home/ubuntu/api/",
			setUp: func(fs afero.Fs, wd string) error {
				if err := fs.Mkdir(filepath.Join(wd, "foo"), directoryPermission); err != nil {
					return microerror.Mask(err)
				}

				if err := fs.Mkdir(filepath.Join(wd, "bar"), directoryPermission); err != nil {
					return microerror.Mask(err)
				}

				if _, err := fs.Create(filepath.Join(wd, "main.go")); err != nil {
					return microerror.Mask(err)
				}

				if _, err := fs.Create(filepath.Join(wd, "foo", "foo.go")); err != nil {
					return microerror.Mask(err)
				}

				if _, err := fs.Create(filepath.Join(wd, "bar", "bar.go")); err != nil {
					return microerror.Mask(err)
				}

				return nil
			},
			expectedDirectories: []string{".", "./bar/...", "./foo/..."},
		},
	}

	for index, test := range tests {
		fs := afero.NewMemMapFs()

		if err := test.setUp(fs, test.workingDirectory); err != nil {
			t.Fatalf("%v: unexpected error during set up: %v", index, err)
		}

		directories, err := NoVendor(fs, test.workingDirectory)
		if err != nil {
			t.Fatalf("%v: unexpected error during no vendor: %v", index, err)
		}

		if !reflect.DeepEqual(directories, test.expectedDirectories) {
			t.Fatalf(
				"%v: returned directories did not match expected.\nexpected:\n%v\nreturned:\n%v\n",
				index,
				test.expectedDirectories,
				directories,
			)
		}
	}
}
