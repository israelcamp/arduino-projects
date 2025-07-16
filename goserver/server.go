package main

import (
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

func capture() {
	resp, err := http.Get("http://192.168.0.107/capture")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Status", resp.Status)
	nowString := getNowFormated()

	currentImagePath = fmt.Sprintf("images/%s.png", nowString)
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

func keepCapturing() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		capture()
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

	go keepCapturing()

	http.HandleFunc("/capture", serveCapture)

	http.ListenAndServe(":8090", nil)
}
