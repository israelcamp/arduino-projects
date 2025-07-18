package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Esp32Cam struct {
		URL             string `yaml:"url"`
		CaptureEndpoint string `yaml:"captureEndpoint"`
		StreamEndpoint  string `yaml:"streamEndpoint"`
	} `yaml:"esp32cam"`
	FileSystem struct {
		ImagesDir string `yaml:"imagesDir"`
	} `yaml:"fileSystem"`
	Capture struct {
		Interval int `yaml:"interval"`
	} `yaml:"capture"`
	RabbitMQ struct {
		Port  string `yaml:"port"`
		User  string `yaml:"user"`
		Pass  string `yaml:"pass"`
		Host  string `yaml:"host"`
		VHost string `yaml:"vhost"`
		Queue string `yaml:"queue"`
	} `yaml:"rabbitmq"`
}

func ReadConfig() Config {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}
	return cfg
}
