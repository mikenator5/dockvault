package main

import (
	"dockVault/helpers"
	"dockVault/storage"
	"log"
)

func main() {
	c, err := helpers.NewConfig()
	if err != nil {
		log.Fatalln("Config failed", err)
	}
	d, err := helpers.NewDocker()
	if err != nil {
		log.Fatal("Failed to start Docker... is it running?")
	}

	s3, err := storage.NewS3(c, d)
	if err != nil {
		log.Fatalln("S3 Failed", err)
	}
	// s3.Upload(storage.UploadParams{
	// 	ImageId:  "alpine",
	// 	BlobName: "alpine_s3",
	// })
	err = s3.Load(storage.LoadParams{
		BlobName: "alpine_s3",
	})
	if err != nil {
		log.Fatal("Failed to download ", err)
	}
}
