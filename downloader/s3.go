package downloader

import (
	"bytes"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Config struct {
	Region       *string
	Bucket       *string
	RevisionPath *string
	Revision     *string
}

type S3Revision struct {
	bucket       string
	revisionPath string
	client       *s3.S3
}

// Get downloads the specified release from S3 bucket, as described in
// the passed in config, and returns the location of the downloaded
// release and an error, if any
func Get(conf *S3Config) (string, error) {
	defaults.DefaultConfig.Region = aws.String(*conf.Region)

	tmpDir := os.TempDir()

	destFile := pathWithFile(&tmpDir, conf.Revision)

	downloadFile, err := os.Create(destFile)
	if err != nil {
		return "", err
	}

	revisionPath := pathWithFile(conf.RevisionPath, conf.Revision)
	downloader := s3manager.NewDownloader(nil)
	_, err = downloader.Download(
		downloadFile,
		&s3.GetObjectInput{
			Bucket: conf.Bucket,
			Key:    &revisionPath,
		})
	if err != nil {
		return "", err
	}

	return destFile, nil
}

// pathWithFile puts together the path to the file with the file name.
func pathWithFile(path *string, revision *string) string {
	var fullRevisionPath bytes.Buffer
	fullRevisionPath.WriteString(*path)
	fullRevisionPath.WriteString(*revision)

	return fullRevisionPath.String()
}
