package internal

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Data       string `json:"data"`
	SlugLength int    `json:"slugLength"`
}

func ParseConfig() (*Config, error) {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	jsonBytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(jsonBytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
