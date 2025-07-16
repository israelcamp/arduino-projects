package main

import (
	"archome/server/config"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var currentImagePath string

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

func keepCapturing(cfg config.Config) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	endpoint := fmt.Sprintf("%s%s", cfg.Esp32Cam.URL, cfg.Esp32Cam.CaptureEndpoint)
	for range ticker.C {
		capture(endpoint, cfg.FileSystem.ImagesDir)
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

	cfg := config.ReadConfig()

	go keepCapturing(cfg)

	http.HandleFunc("/capture", serveCapture)

	http.ListenAndServe(":8090", nil)
}
