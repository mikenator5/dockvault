package helpers

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"

	"github.com/docker/docker/client"
)

type Docker struct {
	ctx context.Context
	cli *client.Client
}

func NewDocker() (Docker, error) {
	d := Docker{}
	d.ctx = context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return Docker{}, err
	}

	if _, err = cli.Info(d.ctx); err != nil {
		return Docker{}, err
	}
	d.cli = cli

	return d, nil
}

func (d Docker) SaveImageInMemory(imageId string) (*bytes.Buffer, error) {
	// Save docker image to buffer
	var tarballBuffer bytes.Buffer
	saveResponse, err := d.cli.ImageSave(d.ctx, []string{imageId})
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

func (d Docker) LoadImageFromBuffer(buffer *bytes.Buffer) error {
	_, err := d.cli.ImageLoad(d.ctx, buffer, false)
	if err != nil {
		return err
	}
	return nil
}
