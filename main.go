package main

import (
	"flag"
	"log"

	"github.com/misham/peter/deployer"
	"github.com/misham/peter/downloader"
	//"github.com/misham/peter/extractor"
)

const (
	defaultCurrentPath = "/var/www/current"
	defaultRegion      = "us-east-1"
)

func main() {
	region := flag.String("region", defaultRegion, "AWS region to use")
	bucket := flag.String("bucket", "", "AWS S3 bucket containing the revision")
	revisionPath := flag.String("revision-path", "/", "AWS S3 path in the bucket where revisions are stored")
	revision := flag.String("revision", "", "AWS S3 key that is the revision to deploy")
	//currentPath := flag.String("current-path", defaultCurrentPath, "Path to the running revision")

	flag.Parse()

	compressedRelease, err := downloader.Get(&downloader.S3Config{
		Region:       region,
		Bucket:       bucket,
		RevisionPath: revisionPath,
		Revision:     revision,
	})
	if err != nil {
		log.Fatal("Failed to download revision", err)
	}

	extractedPath, err := extractor.Extract(compressedRelease)
	if err != nil {
		log.Fatal("Failed to extract release", err)
	}

	//err = deployer.Deploy(extractedPath, currentPath)
}
