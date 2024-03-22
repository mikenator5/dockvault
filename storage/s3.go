package storage

import (
	"bytes"
	"context"
	"dockvault/helpers"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3 represents the S3 storage implementation.
type S3 struct {
	cfg      helpers.Config
	d        helpers.Docker
	s3Client *s3.Client
}

// NewS3 creates a new S3 storage instance using the provided configuration and Docker client.
// It returns the created S3 instance and an error if any occurred during initialization.
func NewS3(cfg helpers.Config, d helpers.Docker) (S3, error) {
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return S3{}, err
	}

	client := s3.NewFromConfig(awsConfig)

	return S3{cfg: cfg, d: d, s3Client: client}, nil
}

// Upload uploads a docker image to S3.
// It takes an UploadParams struct as a parameter, which contains the necessary information for the upload.
// It returns an error if the upload fails.
func (s S3) Upload(params UploadParams) error {
	compressedImg, err := s.d.SaveImageInMemory(params.ImageId)
	if err != nil {
		return err
	}
	fmt.Printf("Uploading %s\n", params.ImageId)

	if len(params.BlobName) == 0 {
		params.BlobName = params.ImageId
	}

	uploader := manager.NewUploader(s.s3Client)
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.cfg.AWS.Bucket),
		Key:    aws.String(params.BlobName),
		Body:   compressedImg,
	})
	if err != nil {
		return err
	}

	fmt.Printf("\nSuccessfully saved %s to %s\n", params.ImageId, result.Location)
	return nil
}

// List retrieves a list of objects from the S3 bucket.
// It prints the key and size of each object found in the bucket.
// If no objects are found, it prints a message indicating that.
func (s S3) List() error {
	output, err := s.s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(s.cfg.AWS.Bucket),
	})
	if err != nil {
		return err
	}

	if len(output.Contents) < 1 {
		fmt.Printf("No objects found in bucket \"%s\"\n", s.cfg.AWS.Bucket)
	}

	for _, object := range output.Contents {
		sizeMB := float64(*object.Size) / (1024 * 1024)
		fmt.Printf("key=%s size=%.2f MB\n", aws.ToString(object.Key), sizeMB)
	}
	return nil
}

// Load downloads an object from an S3 bucket and loads it as a Docker image.
// It takes a LoadParams struct as a parameter, which contains the name of the blob to download.
// It returns an error if there was a problem downloading or loading the image.
func (s S3) Load(params LoadParams) error {
	fmt.Println("Downloading " + params.BlobName)
	headObj, err := s.s3Client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(s.cfg.AWS.Bucket),
		Key:    aws.String(params.BlobName),
	})
	if err != nil {
		return err
	}

	downloader := manager.NewDownloader(s.s3Client)
	data := make([]byte, int(*headObj.ContentLength))
	w := manager.NewWriteAtBuffer(data)

	_, err = downloader.Download(context.TODO(), w, &s3.GetObjectInput{
		Bucket: aws.String(s.cfg.AWS.Bucket),
		Key:    aws.String(params.BlobName),
	})
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(data)
	fmt.Println("Successfully downloaded. Loading image...")
	err = s.d.LoadImageFromBuffer(buffer)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully loaded %s\n", params.BlobName)

	return nil
}
