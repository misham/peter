package extractor

import (
	"io"
	"os"
	"path/filepath"
)

const (
	defaultPermission = 0775
)

type CompressedFileInfo struct {
	path     string
	info     os.FileInfo
	linkName string
}

func (f *CompressedFileInfo) extractArchive(dest *string, input io.Reader) error {
	if f.info.IsDir() {
		err := f.createDir(false)
		if err != nil {
			return err
		}
	} else {
		err := f.createDir(true)
		if err != nil {
			return err
		}

		if f.isSymlink() {
			return os.Symlink(f.linkName, f.path)
		}

		err = f.extractFile(input)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *CompressedFileInfo) createDir(parent bool) error {
	if parent {
		return os.MkdirAll(filepath.Dir(f.path), defaultPermission)
	} else {
		return os.MkdirAll(f.path, f.info.Mode())
	}
}

func (f *CompressedFileInfo) isSymlink() bool {
	return (f.info.Mode() & os.ModeSymlink) != 0
}

func (f *CompressedFileInfo) extractFile(input io.Reader) error {
	fileCopy, err := os.OpenFile(f.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, f.info.Mode())
	if err != nil {
		return err
	}
	defer fileCopy.Close()

	_, err = io.Copy(fileCopy, input)
	if err != nil {
		return err
	}

	return nil
}
