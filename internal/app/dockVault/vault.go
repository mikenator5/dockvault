package dockVault

import (
	"dockVault/internal/pkg/cloudStorage"
	"dockVault/internal/pkg/output"
	"errors"
	"flag"
	"fmt"
	"os"
)

const (
	AzureStorage = "az"
	S3Storage    = "s3"
)

func NewVault() {
	switch os.Args[1] {
	case "upload":
		handleUpload()
	case "list":
		handleList()
	case "load":
		fmt.Println("loading image...")
	default:
		output.PrintUsage()
	}
}

func handleUpload() {
	uploadCmd := flag.NewFlagSet("upload", flag.ExitOnError)
	uploadImageId := uploadCmd.String("image", "", "Docker image id or name:tag")
	uploadBlobName := uploadCmd.String("name", "", "Name of the blob")
	uploadStorage := uploadCmd.String("storage", "", "<az | s3 >")
	uploadCmd.Parse(os.Args[2:])
	storage, err := getStorage(*uploadStorage)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = storage.Upload(storage.GetContainerName(), *uploadImageId, *uploadBlobName)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func handleList() {
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listStorage := listCmd.String("storage", "", "<az | s3>")
	listCmd.Parse(os.Args[2:])
	storage, err := getStorage(*listStorage)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = storage.List()
	if err != nil {
		fmt.Println("Failed to list objects in " + storage.GetContainerName())
	}
}

func getStorage(storage string) (cloudStorage.Storage, error) {
	switch storage {
	case AzureStorage:
		account := os.Getenv("AZ_ACCOUNT")
		containerName := os.Getenv("AZ_CONTAINER_NAME")
		az, err := cloudStorage.NewAzureWithAD(account, containerName)
		if err != nil {
			fmt.Println("Error with Azure Credentials")
			return nil, err
		}
		return &az, err
	default:
		return nil, errors.New("unknown command provided")
	}
}
