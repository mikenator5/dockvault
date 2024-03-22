package storage

import (
	"bytes"
	"context"
	"dockVault/helpers"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3 struct {
	cfg      helpers.Config
	d        helpers.Docker
	s3Client *s3.Client
}

func NewS3(cfg helpers.Config, d helpers.Docker) (S3, error) {
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return S3{}, err
	}

	client := s3.NewFromConfig(awsConfig)

	return S3{cfg: cfg, d: d, s3Client: client}, nil
}

func (s *S3) Upload(params UploadParams) error {
	fmt.Printf("Uploading %s\n", params.ImageId)
	compressedImg, err := s.d.SaveImageInMemory(params.ImageId)
	if err != nil {
		return err
	}

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

func (s *S3) List() error {
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
		fmt.Printf("key=%s size=%d", aws.ToString(object.Key), object.Size)
	}
	return nil
}

func (s *S3) Load(params LoadParams) error {
	fmt.Println("Downloading " + params.BlobName)

	result, err := s.s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.cfg.AWS.Bucket),
		Key:    aws.String(params.BlobName),
	})
	if err != nil {
		return err
	}

	downloader := manager.NewDownloader(s.s3Client)
	data := make([]byte, int(*result.ContentLength))
	w := manager.NewWriteAtBuffer(data)

	_, err = downloader.Download(context.TODO(), w, &s3.GetObjectInput{
		Bucket: aws.String(s.cfg.AWS.Bucket),
		Key:    aws.String(params.BlobName),
	})
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(data)
	err = s.d.LoadImageFromBuffer(buffer)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully loaded %s\n", params.BlobName)

	return nil
}
