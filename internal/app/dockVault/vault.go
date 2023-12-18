package dockVault

import (
	"context"
	"dockVault/internal/pkg/cloudStorage"
	"dockVault/internal/pkg/dockerHelpers"
	"flag"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
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
		size := int64(compressedImage.Len())
		_, err = client.UploadBuffer(context.TODO(), containerName, blobName, compressedImage.Bytes(), &azblob.UploadBufferOptions{
			BlockSize:   int64(1024),
			Concurrency: uint16(3),
			Progress: func(bytesTransferred int64) {
				percentage := int((float64(bytesTransferred) / float64(size)) * 100)
				printProgressBar(percentage)
			},
		})
		if err != nil {
			return err
		}
		fmt.Printf("\nSuccessfully saved %s\n", imageId)
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

func printProgressBar(percentage int) {
	barLength := 50
	numBars := int(float64(barLength) * (float64(percentage) / 100))
	bar := "[" + repeatStr("=", numBars) + repeatStr(" ", barLength-numBars) + "]"
	fmt.Printf("\r%s %d%%", bar, percentage)
}

func repeatStr(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
