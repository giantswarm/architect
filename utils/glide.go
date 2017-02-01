package utils

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// NoVendor replicates glide's novendor command, so we don't have to
// package it at all.
func NoVendor(workingDirectory string) ([]string, error) {
	directories, err := ioutil.ReadDir(workingDirectory)
	if err != nil {
		return nil, err
	}

	testPackages := []string{}

	for _, directory := range directories {
		if !directory.IsDir() {
			continue
		}

		if directory.Name() == "vendor" {
			continue
		}

		files, err := ioutil.ReadDir(directory.Name())
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			if filepath.Ext(file.Name()) == ".go" {
				testPackages = append(testPackages, directory.Name())
				break
			}
		}
	}

	testPackageArguments := []string{"."}
	for _, testPackage := range testPackages {
		testPackageArguments = append(testPackageArguments, fmt.Sprintf("./%v/...", testPackage))
	}

	return testPackageArguments, nil
}
