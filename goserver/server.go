package main

import (
	"archome/server/capture"
	"archome/server/config"
	"archome/server/utils"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	mu      sync.RWMutex
	frame   []byte
	aiframe []byte
)

func keepSavingFrame(cfg config.Config) {
	ticker := time.NewTicker(time.Duration(cfg.Capture.Interval) * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		if cfg.Capture.Save {
			capture.SaveCapture(cfg.FileSystem.ImagesDir, &mu, frame)
		}
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

func streamAIHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
	for {
		mu.RLock()
		f := aiframe
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

func serveFrame(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "image/jpeg")
	b64 := utils.EncodeB64(frame)
	w.Write([]byte(b64))
}

func serveAIFrame(w http.ResponseWriter, req *http.Request) {
	mu.RLock()
	defer mu.RUnlock()
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(aiframe)))
	w.Write(aiframe)
}

func receiveAIFrame(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve the file from the form data
	file, _, err := req.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 3) read all bytes
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "error reading file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 4) store in global (with locking)
	mu.Lock()
	aiframe = data
	mu.Unlock()

	// 5) respond OK
	w.WriteHeader(http.StatusOK)
}

func main() {
	cfg := config.ReadConfig()

	go keepSavingFrame(cfg)
	go capture.FetchFrameLoop(cfg, &mu, &frame)

	http.HandleFunc("/ai", receiveAIFrame)
	http.HandleFunc("/capture", serveFrame)
	http.HandleFunc("/aicapture", serveAIFrame)
	http.HandleFunc("/stream", streamHandler)
	http.HandleFunc(("/streamai"), streamAIHandler)

	http.ListenAndServe(":8090", nil)
}
