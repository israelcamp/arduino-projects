import cv2 as cv
import pika
import srsly

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


if __name__ == "__main__":
    main()
