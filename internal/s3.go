package internal

import (
	"context"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func Upload(ctx context.Context, bucket string, name string, content string) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &name,
		Body:   strings.NewReader(content),
	})
	return err
}

func Download(ctx context.Context, bucket string, name string) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", err
	}

	client := s3.NewFromConfig(cfg)

	out, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &name,
	})
	if err != nil {
		return "", err
	}

	content, err := io.ReadAll(out.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
