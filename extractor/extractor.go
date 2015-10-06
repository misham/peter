package extractor

import (
	"fmt"
	"net/http"
	"os"
)

// Extract extracts file found at the src into dest and
// returns the full path to the extracted directory with any errors.
func Extract(dest *string, src *string) error {
	file, err := os.Open(*src)
	if err != nil {
		return err
	}
	defer file.Close()

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		return err
	}

	filetype := http.DetectContentType(buff)

	switch filetype {
	case "application/zip":
		err = unzip(src, dest)
	case "application/x-gzip":
		err = untgz(src, dest)
	default:
		err = fmt.Errorf("Unknown file type: %s", filetype)
	}

	return err
}
