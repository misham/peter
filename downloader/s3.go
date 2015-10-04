package downloader

import (
	"os"
	"path/filepath"

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

	destFile := filepath.Join(tmpDir, *conf.Revision)

	downloadFile, err := os.Create(destFile)
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(*conf.RevisionPath, *conf.Revision)
	downloader := s3manager.NewDownloader(nil)
	_, err = downloader.Download(
		downloadFile,
		&s3.GetObjectInput{
			Bucket: conf.Bucket,
			Key:    &filePath,
		})
	if err != nil {
		return "", err
	}

	return destFile, nil
}
