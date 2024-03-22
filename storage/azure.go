package storage

import (
	"bytes"
	"context"
	"dockvault/helpers"
	"dockvault/output"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

// Azure represents the Azure storage implementation.
type Azure struct {
	cfg        helpers.Config
	d          helpers.Docker
	blobClient *azblob.Client
}

// NewAzureWithAD creates a new Azure storage instance with Active Directory authentication.
// It takes a configuration object `cfg` and a Docker helper object `d` as parameters.
// It returns an instance of the Azure struct and an error, if any.
func NewAzureWithAD(cfg helpers.Config, d helpers.Docker) (Azure, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return Azure{}, err
	}

	serviceUrl := fmt.Sprintf("https://%s.blob.core.windows.net/", cfg.Azure.StorageAccount)
	client, err := azblob.NewClient(serviceUrl, cred, nil)
	if err != nil {
		return Azure{}, err
	}

	return Azure{cfg: cfg, d: d, blobClient: client}, nil
}

// Upload uploads the image specified by the given UploadParams to Azure Blob Storage.
// It compresses the image, assigns a BlobName if not provided, and uploads the compressed image to the specified container.
// The method returns an error if any error occurs during the upload process.
func (az Azure) Upload(params UploadParams) error {
	compressedImg, err := az.d.SaveImageInMemory(params.ImageId)
	if err != nil {
		return err
	}
	fmt.Printf("Uploading %s\n", params.ImageId)

	if len(params.BlobName) == 0 {
		params.BlobName = params.ImageId
	}

	size := int64(compressedImg.Len())
	_, err = az.blobClient.UploadBuffer(context.TODO(), az.cfg.Azure.Container, params.BlobName, compressedImg.Bytes(), &azblob.UploadBufferOptions{
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

	fmt.Printf("\nSuccessfully saved %s\n", params.ImageId)
	return nil
}

// List retrieves a list of blobs from the Azure storage container.
// It returns an error if there was a problem listing the blobs.
func (az Azure) List() error {
	pager := az.blobClient.NewListBlobsFlatPager(az.cfg.Azure.Container, nil)
	blobCount := 0
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			return err
		}

		for _, blob := range page.Segment.BlobItems {
			fmt.Println(*blob.Name)
			blobCount++
		}
	}

	if blobCount <= 0 {
		fmt.Printf("No blobs found in container \"%s\"\n", az.cfg.Azure.Container)
	}

	return nil
}

// Load downloads a blob from Azure storage and loads it into the storage.
func (az Azure) Load(params LoadParams) error {
	fmt.Printf("Downloading %s\n", params.BlobName)
	get, err := az.blobClient.DownloadStream(context.TODO(), az.cfg.Azure.Container, params.BlobName, nil)
	if err != nil {
		return err
	}

	data := bytes.Buffer{}
	retryReader := get.NewRetryReader(context.TODO(), nil)
	_, err = data.ReadFrom(retryReader)
	if err != nil {
		return err
	}

	err = az.d.LoadImageFromBuffer(&data)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully loaded %s\n", params.BlobName)
	return nil
}
