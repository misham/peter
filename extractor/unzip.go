package extractor

import (
	"archive/zip"
	"io/ioutil"
	"path"
)

func unzip(src *string, dest *string) error {
	reader, err := zip.OpenReader(*src)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		zipReader, err := file.Open()
		if err != nil {
			return err
		}
		defer zipReader.Close()

		info := CompressedFileInfo{
			path: path.Join(*dest, file.Name),
			info: file.FileInfo(),
		}

		if info.isSymlink() {
			linkName, err := ioutil.ReadAll(zipReader)
			if err != nil {
				return err
			}
			info.linkName = string(linkName)
		}

		err = info.extractArchive(dest, zipReader)
		if err != nil {
			return err
		}
	}

	return nil
}
