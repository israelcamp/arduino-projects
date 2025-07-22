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

type Paths struct {
	Date string
	Time string
}

func getNowFormated() Paths {
	now := time.Now()
	dateString := fmt.Sprintf("%d%02d%02d", now.Year(), now.Month(), now.Day())
	timeString := fmt.Sprintf("%02d%02d%02d", now.Hour(), now.Minute(), now.Second())
	paths := Paths{
		Date: dateString,
		Time: timeString,
	}
	return paths
}

func Capture(imagesDir string, mu *sync.RWMutex, frame []byte) {
	nowPaths := getNowFormated()

	mu.RLock()
	content := frame
	mu.RUnlock()

	imageDir := fmt.Sprintf("%s/%s", imagesDir, nowPaths.Date)
	os.MkdirAll(imageDir, os.ModePerm)

	imagePath := fmt.Sprintf("%s/%s.png", imageDir, nowPaths.Time)

	out, err := os.Create(imagePath)
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
