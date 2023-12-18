package cloudStorage

type Storage interface {
	GetContainerName() string
	Upload(containerName string, imageId string, blobName string) error
	List() error
	Load()
}
