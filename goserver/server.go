package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Esp32Cam struct {
		URL             string `yaml:"url"`
		CaptureEndpoint string `yaml:"capture"`
	} `yaml:"esp32cam"`
	FileSystem struct {
		ImagesDir string `yaml:"imagesDir"`
	} `yaml:"fileSystem"`
}

var currentImagePath string

func readConfig() Config {
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

func getNowFormated() string {
	now := time.Now()
	return fmt.Sprintf("%d-%d-%d_%d-%d-%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}

func capture(endpoint string, imagesDir string) {
	resp, err := http.Get(endpoint)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Status", resp.Status)
	nowString := getNowFormated()

	currentImagePath = fmt.Sprintf("%s/%s.png", imagesDir, nowString)
	out, err := os.Create(currentImagePath)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("Image saved")
}

func keepCapturing(cfg Config) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		capture(fmt.Sprintf("%s%s", cfg.Esp32Cam.URL, cfg.Esp32Cam.CaptureEndpoint), cfg.FileSystem.ImagesDir)
	}
}

func serveCapture(w http.ResponseWriter, req *http.Request) {
	buf, err := os.ReadFile(currentImagePath)
	if err != nil {
		http.Error(w, "ERROR: NOT FOUND", 404)
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write(buf)
}

func main() {

	config := readConfig()

	go keepCapturing(config)

	http.HandleFunc("/capture", serveCapture)

	http.ListenAndServe(":8090", nil)
}
