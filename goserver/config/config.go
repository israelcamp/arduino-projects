package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type EnvVar string

func (e *EnvVar) UnmarshalYAML(node *yaml.Node) error {
	if node.Tag == "!env_var" {
		// Expect two child nodes: key and optional default
		if len(node.Content) < 1 {
			return fmt.Errorf("!env_var requires at least a key")
		}
		key := node.Content[0].Value

		var def string
		if len(node.Content) >= 2 {
			def = node.Content[1].Value
		}

		if val := os.Getenv(key); val != "" {
			*e = EnvVar(val)
		} else {
			*e = EnvVar(def)
		}
		return nil
	}

	// Fallback: decode normally (e.g. plain scalar)
	var plain string
	if err := node.Decode(&plain); err != nil {
		return err
	}
	*e = EnvVar(plain)
	return nil
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
