package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

func envOr(key, def string) string {
	if v, ok := os.LookupEnv(key); ok { // set (could be "")
		return v
	}
	return def
}

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
		Interval     int  `yaml:"interval"`
		Save         bool `yaml:"save"`
		SaveAI       bool `yaml:"saveAI"`
		SaveOnPerson bool `yaml:"saveOnPerson"`
	} `yaml:"capture"`
}

func ReadConfig() Config {
	configPath := envOr("GS_CONFIG", "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}
	return cfg
}
