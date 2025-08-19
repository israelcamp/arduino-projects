# Arduino Projects

## Surveillance System

To build my own security system I am using as hardware:
1. ESP32-CAM
1. ESP32 WiFi
1. TFT Display

The software I am using are:
1. Go server
1. Python for AI processes

The idea is to install the ESP32-CAM (currently I am only using one, but will expand) at some spot in my house that allows me to see the front gate.
The ESP32-CAM is connected to my house WiFi.
I altered the default image quality and size to stream a 640x480 image.

The Go severr will be my main form of connection to the camera feed. Because the ESP32-CAM cannot handle well multiple connections I decided to make a Go server that will stay connected to the camera feeds and keep the last frame on memory. This will allow me to fetch the image from the ESP32-CAM only once and feed it to as many connections as I need through the go server. This server will have 4 main endpoints:
1. `/stream` - opens a connection to keep feeding the client with the latest frame from the ESP32-CAM
1. `/capture` - serves the last frame in base64
1. `/streamai` - the same as `/stream` with the difference that the image fed has been processed by a YOLO model for object detection
1. `/aicapture` - the same as `/capture` but with the image processed by AI

The AI model will run in a python script. I found easier to use python as it has more support for different models and is easier to use them.
The script will keep making requests to `/capture`, processing the imagem, drawing the bounding boxes around the detected objects and uploading the new image to another endpoint.
The uploaded image will be the one served on `/aicapture` and `/streamai`.

The model I am currently using is a YOLO-v4-tiny that I found the weights in github. It is super easy to use with `opencv2` and has good perfomance and speed.

