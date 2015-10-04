package extractor

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

func untgz(src *string, dest *string) error {
	fd, err := os.Open(*src)
	if err != nil {
		return err
	}
	defer fd.Close()

	gReader, err := gzip.NewReader(fd)
	if err != nil {
		return err
	}
	defer gReader.Close()

	tarReader := tar.NewReader(gReader)

	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if hdr.Name == "." {
			continue
		}

		info := CompressedFileInfo{
			path:     filepath.Join(*dest, hdr.Name),
			info:     hdr.FileInfo(),
			linkName: hdr.Linkname,
		}
		err = info.extractArchive(dest, tarReader)
		if err != nil {
			return err
		}
	}

	return nil
}
