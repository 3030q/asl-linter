package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Minio struct {
	Endpoint string `json:"endpoint"`
	Username string `json:"username"`
	Password string `json:"password"`
	UseSsl   bool   `json:"useSsl"`
}

type StateLintConfig struct {
	MinioConfig Minio `json:"minio"`
}

const DefaultConfigPath = "./statelintConfig.json"

func ReadConfig() (*StateLintConfig, error) {
	config := StateLintConfig{}

	fileData, err := ioutil.ReadFile(DefaultConfigPath)
	if err != nil {
		return nil, fmt.Errorf("config file read error %w", err)
	}

	err = json.Unmarshal(fileData, &config)
	if err != nil {
		return nil, fmt.Errorf("config file unmarshal error %w", err)
	}

	return &config, nil
}
