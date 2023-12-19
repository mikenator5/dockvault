package cloudStorage

import (
	"bytes"
	"context"
	"dockVault/internal/pkg/dockerHelpers"
	"dockVault/internal/pkg/output"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type Azure struct {
	ContainerName string
	Cred          *azidentity.DefaultAzureCredential
	Client        *azblob.Client
}

func NewAzureWithAD(storageUrl string, containerName string) (Azure, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return Azure{}, err
	}

	client, err := azblob.NewClient(storageUrl, cred, nil)
	if err != nil {
		return Azure{}, err
	}
	return Azure{
		ContainerName: containerName,
		Cred:          cred,
		Client:        client,
	}, nil
}

func (az *Azure) GetContainerName() string {
	return az.ContainerName
}

func (az *Azure) Upload(containerName string, imageId string, blobName string) error {
	fmt.Printf("Uploading %s\n", imageId)
	compressedImage, err := dockerHelpers.SaveImageInMemory(imageId)
	if err != nil {
		return err
	}
	if len(blobName) == 0 {
		blobName = imageId
	}

	client := az.Client

	size := int64(compressedImage.Len())
	_, err = client.UploadBuffer(context.TODO(), containerName, blobName, compressedImage.Bytes(), &azblob.UploadBufferOptions{
		BlockSize:   int64(1024),
		Concurrency: uint16(3),
		Progress: func(bytesTransferred int64) {
			percentage := int((float64(bytesTransferred) / float64(size)) * 100)
			output.PrintProgressBar(percentage)
		},
	})
	if err != nil {
		return err
	}
	fmt.Printf("\nSuccessfully saved %s\n", imageId)
	return nil
}

func (az *Azure) List() error {
	pager := az.Client.NewListBlobsFlatPager(az.GetContainerName(), nil)

	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			return err
		}
		for _, blob := range page.Segment.BlobItems {
			fmt.Println(*blob.Name)
		}
	}
	return nil
}

func (az *Azure) Load(blobName string) error {
	fmt.Println("Downloading " + blobName)
	get, err := az.Client.DownloadStream(context.TODO(), az.GetContainerName(), blobName, nil)
	if err != nil {
		return err
	}
	data := bytes.Buffer{}
	retryReader := get.NewRetryReader(context.TODO(), nil)
	_, err = data.ReadFrom(retryReader)
	if err != nil {
		return err
	}
	err = retryReader.Close()
	if err != nil {
		return err
	}

	err = dockerHelpers.LoadImageFromBuffer(&data)
	if err != nil {
		return err
	}
	fmt.Printf("Successfully loaded %s\n", blobName)
	return nil
}
