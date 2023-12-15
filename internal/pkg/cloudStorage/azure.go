package cloudStorage

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func NewActiveDirClient(storageUrl string) (*azblob.Client, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	client, err := azblob.NewClient(storageUrl, cred, nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}
