package cloudStorage

import (
	"bytes"
	"context"
	"dockVault/internal/pkg/dockerHelpers"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
)

type S3 struct {
	ContainerName string
	Config        aws.Config
	Bucket        string
}

func NewS3(bucket string, containerName string) (S3, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return S3{}, err
	}

	return S3{ContainerName: containerName, Config: cfg, Bucket: bucket}, nil
}

func (s S3) GetContainerName() string {
	return s.ContainerName
}

func (s S3) Upload(containerName string, imageId string, blobName string) error {
	fmt.Printf("Uploading %s\n", imageId)
	compressedImage, err := dockerHelpers.SaveImageInMemory(imageId)
	if err != nil {
		return err
	}
	if len(blobName) == 0 {
		blobName = imageId
	}

	client := s3.NewFromConfig(s.Config)

	uploader := manager.NewUploader(client)
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(blobName),
		Body:   compressedImage,
	})
	if err != nil {
		return err
	}

	fmt.Printf("\nSuccessfully saved %s to %s\n", imageId, result.Location)
	return nil
}

func (s S3) List() error {
	client := s3.NewFromConfig(s.Config)
	output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(s.Bucket),
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, object := range output.Contents {
		log.Printf("key=%s size=%d", aws.ToString(object.Key), object.Size)
	}
	return nil

}

func (s S3) Load(blobName string) error {
	fmt.Println("Downloading " + blobName)
	client := s3.NewFromConfig(s.Config)

	result, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(blobName),
	})
	if err != nil {
		return err
	}

	downloader := manager.NewDownloader(client)
	data := make([]byte, int(*result.ContentLength))
	w := manager.NewWriteAtBuffer(data)

	_, err = downloader.Download(context.TODO(), w, &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(blobName),
	})
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(data)
	err = dockerHelpers.LoadImageFromBuffer(buffer)
	if err != nil {
		return err
	}
	fmt.Printf("Successfully loaded %s\n", blobName)

	return nil
}
