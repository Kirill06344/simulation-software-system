package util

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Lambda        int32 `yaml:"lambda"`
	SourcesAmount int32 `yaml:"sourcesAmount"`
	RequestAmount int32 `yaml:"requestAmount"`
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
