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
	installPath := flag.String("install-path", "/var/www", "Path to deploy revisions to")

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 {
		log.Println("Argument required")
		usage()
		os.Exit(1)
	}

	revision := flag.Arg(0)

	compressedReleasePath, err := downloader.Get(&downloader.S3Config{
		Region:       region,
		Bucket:       bucket,
		RevisionPath: revisionPath,
		Revision:     &revision,
	})
	defer os.Remove(compressedReleasePath)
	if err != nil {
		log.Fatal("Failed to download revision: ", err)
	}

	err = extractor.Extract(installPath, &compressedReleasePath)
	if err != nil {
		log.Fatal("Failed to extract release: ", err)
	}

	log.Println(installPath)
}

func usage() {
	io.WriteString(os.Stderr, usageText)
	flag.PrintDefaults()
}

const usageText = `peter [options] <revision ID>
Revision ID is the name of the package to deploy.

`
