package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

func LoadConfig(path string) (Config, error) {
	var config Config
	file, err := os.Open(path)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return config, err
	}
	return config, nil
}
