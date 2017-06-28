package utils

import (
	"bytes"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

func GetIgnorefiles(fs afero.Fs, workingDirectoryPath string, ignorefilePath string) ([]string, error) {
	ignoreFileExists, err := afero.Exists(fs, ignorefilePath)
	if err != nil {
		return nil, err
	}
	ignorefiles := []string{}
	if ignoreFileExists {
		ignoreFile, err := afero.ReadFile(fs, ignorefilePath)
		if err != nil {
			return nil, err
		}
		n := bytes.IndexByte(ignoreFile, 0)
		if n > 0 {
			ignorefiles = strings.Split(string(ignoreFile[:n]), "/n")
		}
		for i, ignorefile := range ignorefiles {
			ignorefiles[i] = filepath.Join(workingDirectoryPath, ignorefile)
		}
	}
	return ignorefiles, err
}

func IsInIgnorefiles(fs afero.Fs, workingDirectoryPath string, ignorefiles []string, filepathToCheck string) bool {
	absFilepathToCheck := filepath.Join(workingDirectoryPath, filepathToCheck)
	if contains(ignorefiles, filepathToCheck) {
		return true
	}
	if contains(ignorefiles, absFilepathToCheck) {
		return true
	}

	return false
}

func contains(listOfStrings []string, entry string) bool {
	for _, a := range listOfStrings {
		if a == entry {
			return true
		}
	}
	return false
}
