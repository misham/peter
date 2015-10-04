package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/misham/peter/downloader"
	"github.com/misham/peter/extractor"
)

const (
	defaultCurrentPath = "/var/www/current"
	defaultRegion      = "us-east-1"
)

func main() {
	region := flag.String("region", defaultRegion, "AWS region to use")
	bucket := flag.String("bucket", "", "AWS S3 bucket containing the revision")
	revisionPath := flag.String("revision-path", "/", "AWS S3 path in the bucket where revisions are stored")
	// TODO move this to an argument
	revision := flag.String("revision", "", "AWS S3 key that is the revision to deploy")
	installPath := flag.String("install-path", "/var/www", "Path to deploy revisions to")

	flag.Usage = usage
	flag.Parse()

	compressedReleasePath, err := downloader.Get(&downloader.S3Config{
		Region:       region,
		Bucket:       bucket,
		RevisionPath: revisionPath,
		Revision:     revision,
	})
	if err != nil {
		log.Fatal("Failed to download revision: ", err)
	}

	extractedPath, err := extractor.Extract(installPath, &compressedReleasePath)
	if err != nil {
		log.Fatal("Failed to extract release: ", err)
	}

	log.Println(extractedPath)

	// TODO remove the downloaded file, it's no longer needed
}

func usage() {
	io.WriteString(os.Stderr, usageText)
	flag.PrintDefaults()
}

const usageText = `peter [options] <revision ID>
Revision ID is the name of the package to deploy.

`
