package internal

import (
	"os"
	"strconv"
)

type Config struct {
	SlugLength int
	BucketName string
}

func ParseConfig() (*Config, error) {
	slugLength, err := strconv.Atoi(os.Getenv("SLUG_LENGTH"))
	if err != nil {
		return nil, err
	}

	return &Config{
		SlugLength: slugLength,
		BucketName: os.Getenv("BUCKET_NAME"),
	}, nil
}
