package storage

// UploadParams represents the parameters required for uploading an image.
type UploadParams struct {
	ImageId  string // The unique identifier of the image.
	BlobName string // The name of the blob associated with the image.
}

// LoadParams represents the parameters for loading a blob.
type LoadParams struct {
	BlobName string // The name of the blob to load.
}

// Storage represents a storage interface for uploading, listing, and loading data.
type Storage interface {
	Upload(UploadParams) error // Upload prepares and uploads a Docker image to the cloud storage provider.
	List() error               // List lists all objects stored in the cloud storage provider.
	Load(LoadParams) error     // Load loads the Docker image from the cloud storage provider into Docker.
}
