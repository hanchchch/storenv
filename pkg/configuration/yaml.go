package configuration

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Storage struct {
		AWS struct {
			Bucket string
			Prefix string
			Region string
		}
	}
	Secrets []string
}

func NewConfiguration(filename string) (*Configuration, error) {
	config := &Configuration{}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}
