package config

import (
	"io/ioutil"
	"omnihr-coding-test/pkg/models"

	"gopkg.in/yaml.v2"
)

func LoadConfig(filePath string) (*models.Config, error) {
	config := &models.Config{}
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
