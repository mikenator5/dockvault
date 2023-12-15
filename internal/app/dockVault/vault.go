package dockVault

import (
	"context"
	"dockVault/internal/pkg/cloudStorage"
	"dockVault/internal/pkg/dockerHelpers"
	"flag"
	"fmt"
	"os"
)

const (
	AzureStorage = "az"
	S3Storage    = "s3"
)

func NewVault(args []string) {
	uploadCmd := flag.NewFlagSet("upload", flag.ExitOnError)
	uploadImageId := uploadCmd.String("image", "", "Docker image id or name:tag")
	uploadBlobName := uploadCmd.String("name", "", "Name of the blob")
	uploadStorage := uploadCmd.String("storage", "", "<az | s3 >")
	switch args[0] {
	case "upload":
		uploadCmd.Parse(args[1:])
		if *uploadStorage != AzureStorage && *uploadStorage != S3Storage {
			uploadCmd.Usage()
			return
		}
		err := upload(*uploadImageId, *uploadStorage, *uploadBlobName)
		if err != nil {
			uploadCmd.Usage()
			return
		}
	case "list":
		fmt.Println("listing saved images...")
	case "load":
		fmt.Println("loading image...")
	default:
		printUsage()
	}
}

func upload(imageId string, storage string, blobName string) error {
	fmt.Printf("Uploading %s to %s\n", imageId, storage)
	compressedImage, err := dockerHelpers.SaveImageInMemory(imageId)
	if err != nil {
		return err
	}
	if len(blobName) == 0 {
		blobName = imageId
	}
	switch storage {
	case AzureStorage:
		account := os.Getenv("AZ_ACCOUNT")
		containerName := os.Getenv("AZ_CONTAINER_NAME")
		client, err := cloudStorage.NewActiveDirClient(account)
		if err != nil {
			return err
		}
		_, err = client.UploadBuffer(context.TODO(), containerName, blobName, compressedImage.Bytes(), nil)
		if err != nil {
			return err
		}
		fmt.Printf("Successfully saved %s to %s\n", imageId, account)
	}
	return nil
}

func list() {

}

func load() {

}

func printUsage() {
	fmt.Println("Usage: dockerHelpers <upload | list | load>")
}
