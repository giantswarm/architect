package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/afero"

	"github.com/giantswarm/microerror"
)

// NoVendor replicates glide's novendor command, so we don't have to
// package it at all.
func NoVendor(fs afero.Fs, workingDirectory string) ([]string, error) {
	packages := map[string]struct{}{}

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return microerror.Mask(err)
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".go" {
			shortPath := strings.TrimPrefix(path, workingDirectory)
			parts := strings.FieldsFunc(shortPath, func(r rune) bool {
				return r == filepath.Separator
			})

			var packageName string
			if len(parts) <= 1 { // e.g: main.go
				packageName = "."
			} else {
				if parts[0] == "vendor" || strings.HasPrefix(parts[0], "_") {
					return nil
				}
				packageName = fmt.Sprintf("./%v/...", parts[0])
			}

			if _, exists := packages[packageName]; !exists {
				packages[packageName] = struct{}{}
			}
		}

		return nil
	}

	if err := afero.Walk(fs, workingDirectory, walkFunc); err != nil && err != filepath.SkipDir {
		return nil, microerror.Mask(err)
	}

	packageNames := []string{}
	for packageName, _ := range packages {
		packageNames = append(packageNames, packageName)
	}

	sort.Strings(packageNames)

	return packageNames, nil
}
