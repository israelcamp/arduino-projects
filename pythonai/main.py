import traceback
import time

import cv2 as cv
import pika
import srsly
import requests

from yolo.yoloclass import ObjectDetection


model = ObjectDetection(
    classes_path="yolo/classes.txt",
    config_path="yolo/yolov4-tiny.cfg",
    weights_path="yolo/yolov4-tiny.weights",
)


global_frame = None


def callback(channel, method_frame, header_frame, body):
    image_b64 = body.decode()
    drawn_frame = model.run(image_b64)
    # cv.imshow("frame", drawn_frame)
    # key = cv.waitKey(1)
    # if key == ord("q"):
    #     channel.stop_consuming()
    channel.basic_ack(delivery_tag=method_frame.delivery_tag)
    return drawn_frame


def main():
    config = srsly.read_yaml("config.yaml")["rabbitmq"]

    credentials = pika.PlainCredentials(config["user"], config["pass"])
    parameters = pika.ConnectionParameters(
        config["host"], int(config["port"]), config["vhost"], credentials
    )

    connection = pika.BlockingConnection(parameters)
    channel = connection.channel()

    channel.basic_consume(config["queue"], callback)

    try:
        print("Starting consuming...")
        channel.start_consuming()
    except KeyboardInterrupt:
        channel.stop_consuming()

    # cv.destroyAllWindows()
    connection.close()


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
