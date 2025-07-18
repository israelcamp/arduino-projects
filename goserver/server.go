package main

import (
	"archome/server/capture"
	"archome/server/config"
	"archome/server/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
	"net/http"
	"sync"
	"time"
)

var (
	mu    sync.RWMutex
	frame []byte
)

func keepCapturing(cfg config.Config, ch *amqp.Channel, q amqp.Queue) {
	ticker := time.NewTicker(time.Duration(cfg.Capture.Interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		capture.Capture(cfg.FileSystem.ImagesDir, &mu, frame)
		rabbitmq.PlubishToQueue(ch, q, frame)
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

	conn := rabbitmq.CreateConnection(cfg)
	defer conn.Close()

	ch := rabbitmq.OpenChannel(conn)
	defer ch.Close()

	queue := rabbitmq.OpenQueue(ch)

	go keepCapturing(cfg, ch, queue)
	go capture.FetchLoop(cfg, &mu, &frame)

	http.HandleFunc("/stream", streamHandler)

	http.ListenAndServe(":8090", nil)
}
