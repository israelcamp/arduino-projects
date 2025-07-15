package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

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

	out, err := os.Create(fmt.Sprintf("%s.png", nowString))
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

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func main() {

	capture()

	http.HandleFunc("/hello", hello)

	http.ListenAndServe(":8090", nil)
}
