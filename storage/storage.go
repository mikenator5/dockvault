package storage

type UploadParams struct {
	ImageId  string
	BlobName string
}

type LoadParams struct {
	BlobName string
}

type Storage interface {
	Upload(UploadParams) error
	List() error
	Load(LoadParams) error
}
