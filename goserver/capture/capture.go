package capture

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

func getNowFormated() string {
	now := time.Now()
	return fmt.Sprintf("%d-%d-%d_%d-%d-%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}

func Capture(imagesDir string, mu *sync.RWMutex, frame []byte) {
	nowString := getNowFormated()

	mu.RLock()
	content := frame
	mu.RUnlock()

	currentImagePath := fmt.Sprintf("%s/%s.png", imagesDir, nowString)
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
func FetchLoop(cfg config.Config, mu *sync.RWMutex, frame *[]byte) {
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
			*frame = buf.Bytes()
			mu.Unlock()
		}
		resp.Body.Close()
	}
}
