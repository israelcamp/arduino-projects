import traceback
import time

import cv2 as cv
import srsly
import requests

from yolo.yoloclass import ObjectDetection


model = ObjectDetection(
    classes_path="yolo/classes.txt",
    config_path="yolo/yolov4-tiny.cfg",
    weights_path="yolo/yolov4-tiny.weights",
)


def aistream():
    config = srsly.read_yaml("config.yaml")
    go_config = config["goserver"]

    base_url = go_config["base_url"]
    capture_endpoint = f"{base_url}/{go_config['capture_endpoint']}"
    upload_endpoint = f"{base_url}/{go_config['upload_endpoint']}"

    delay = int(config["ai"]["interval"])

    while True:
        try:
            req = requests.get(capture_endpoint)
            aiframe = model.run(req.text)
            _, encoded_image = cv.imencode(".jpg", aiframe)
            resp = requests.post(
                upload_endpoint,
                files={"image": ("frame.jpg", encoded_image, "image/jpeg")},
            )
            if resp.status_code != 200:
                raise Exception(f"Error sending image {resp.status_code}: {resp.text}")
            time.sleep(delay)
        except KeyboardInterrupt:
            print("ENDING SERVICE")
            return
        except Exception:
            print("ERROR SENDING AIFRAME")
            print(traceback.format_exc())
            time.sleep(60)
            continue


if __name__ == "__main__":
    aistream()
