package main

import (
	"archome/server/config"
	"bytes"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var currentImagePath string
var (
	mu    sync.RWMutex
	frame []byte
)

func getNowFormated() string {
	now := time.Now()
	return fmt.Sprintf("%d-%d-%d_%d-%d-%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}

func capture(imagesDir string) {
	nowString := getNowFormated()

	mu.RLock()
	content := frame
	mu.RUnlock()

	currentImagePath = fmt.Sprintf("%s/%s.png", imagesDir, nowString)
	out, err := os.Create(currentImagePath)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	_, err = io.Copy(out, bytes.NewReader(content))
	if err != nil {
		panic(err)
	}
	fmt.Println("Image saved")
}

func keepCapturing(cfg config.Config) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		capture(cfg.FileSystem.ImagesDir)
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

func fetchLoop(cfg config.Config) {
	for {
		resp, err := http.Get(fmt.Sprintf("%s:81/%s", cfg.Esp32Cam.URL, cfg.Esp32Cam.StreamEndpoint))
		if err != nil {
			log.Println("STREAM ERROR: ", err)
			time.Sleep(10 * time.Second)
			continue
		}

		// Parse the Content-Type header to get the boundary parameter
		mediaType, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
		if err != nil {
			log.Println("Invalid Content-Type:", err)
			resp.Body.Close()
			time.Sleep(time.Second)
			continue
		}
		if !strings.HasPrefix(mediaType, "multipart/") {
			log.Println("Not a multipart stream:", mediaType)
			resp.Body.Close()
			time.Sleep(time.Second)
			continue
		}
		boundary := params["boundary"] // no leading “--” here

		r := multipart.NewReader(resp.Body, boundary)
		for {
			p, err := r.NextPart()
			if err != nil {
				log.Println("ERROR FETCHING NEXT FRAME: ", err)
				break
			}
			buf := new(bytes.Buffer)
			io.Copy(buf, p)
			mu.Lock()
			frame = buf.Bytes()
			mu.Unlock()
		}
		resp.Body.Close()
	}
}

func streamHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
	for {
		mu.RLock()
		f := frame
		mu.RUnlock()
		if len(f) > 0 {
			w.Write([]byte("--frame\r\n"))
			w.Write([]byte("Content-Type: image/jpeg\r\n\r\n"))
			w.Write(f)
			w.Write([]byte("\r\n"))
		}
		if fl, ok := w.(http.Flusher); ok {
			fl.Flush()
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {

	cfg := config.ReadConfig()

	go keepCapturing(cfg)
	go fetchLoop(cfg)

	http.HandleFunc("/stream", streamHandler)

	http.ListenAndServe(":8090", nil)
}
