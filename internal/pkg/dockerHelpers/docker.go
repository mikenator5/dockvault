package dockerHelpers

import (
	"bytes"
	"compress/gzip"
	"context"
	"github.com/docker/docker/client"
	"io"
)

func Initialize() (context.Context, *client.Client, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	return ctx, cli, err
}

func SaveImageInMemory(imageId string) (*bytes.Buffer, error) {
	ctx, cli, err := Initialize()
	if err != nil {
		return nil, err
	}

	// Save docker image to buffer
	var tarballBuffer bytes.Buffer
	saveResponse, err := cli.ImageSave(ctx, []string{imageId})
	if err != nil {
		return nil, err
	}
	defer saveResponse.Close()
	_, err = io.Copy(&tarballBuffer, saveResponse)
	if err != nil {
		return nil, err
	}

	// Compress image
	var compressedBuffer bytes.Buffer
	writer := gzip.NewWriter(&compressedBuffer)
	defer writer.Close()
	_, err = io.Copy(&compressedBuffer, &tarballBuffer)
	if err != nil {
		return nil, err
	}

	return &compressedBuffer, nil
}

func LoadImageFromBuffer(buffer *bytes.Buffer) error {
	ctx, cli, err := Initialize()
	_, err = cli.ImageLoad(ctx, buffer, false)
	if err != nil {
		return err
	}
	return nil
}
