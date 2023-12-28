package util

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	BufferCapacity int32 `yaml:"bufferCapacity"`
}

func LoadConfig(path string) (Config, error) {
	var cfg Config

	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
