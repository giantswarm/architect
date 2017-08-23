package utils

import (
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/giantswarm/microerror"
)

// Originally stolen from https://gist.github.com/m4ng0squ4sh/92462b38df26839a3ca324697c8cba04

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func CopyFile(fs afero.Fs, src, dst string) error {
	in, err := fs.Open(src)
	if err != nil {
		return microerror.Mask(err)
	}
	defer in.Close()

	out, err := fs.Create(dst)
	if err != nil {
		return microerror.Mask(err)
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return microerror.Mask(err)
	}

	err = out.Sync()
	if err != nil {
		return microerror.Mask(err)
	}

	si, err := fs.Stat(src)
	if err != nil {
		return microerror.Mask(err)
	}
	err = fs.Chmod(dst, si.Mode())
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(fs afero.Fs, src, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := fs.Stat(src)
	if err != nil {
		return microerror.Mask(err)
	}
	if !si.IsDir() {
		return microerror.Mask(sourceNotDirectoryError)
	}

	_, err = fs.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return microerror.Mask(err)
	}
	if err == nil {
		return microerror.Mask(destinationExistsError)
	}

	err = fs.MkdirAll(dst, si.Mode())
	if err != nil {
		return microerror.Mask(err)
	}

	entries, err := afero.ReadDir(fs, src)
	if err != nil {
		return microerror.Mask(err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(fs, srcPath, dstPath)
			if err != nil {
				return microerror.Mask(err)
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(fs, srcPath, dstPath)
			if err != nil {
				return microerror.Mask(err)
			}
		}
	}

	return nil
}
