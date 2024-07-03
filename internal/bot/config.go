package bot

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Bot struct {
		ApiKey string `yaml:"apiKey"`
	} `yaml:"bot"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	}
}

func NewConfig(configPath string) (*Config, error) {

	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}
