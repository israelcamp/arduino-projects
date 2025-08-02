import base64
import io

from PIL import Image
import numpy as np
import cv2 as cv

Conf_threshold = 0.4
NMS_threshold = 0.4
COLORS = [
    (0, 255, 0),
    (0, 0, 255),
    (255, 0, 0),
    (255, 255, 0),
    (255, 0, 255),
    (0, 255, 255),
]


class ODOutput:
    def __init__(self, classes, scores, boxes):
        self.classes = classes
        self.scores = scores
        self.boxes = boxes


class ObjectDetection:
    def __init__(self, classes_path: str, config_path: str, weights_path: str):
        self.classes = self.read_classes(classes_path)
        self.model = self.init_model(config_path, weights_path)

    def read_classes(self, path: str) -> list[str]:
        with open(path, "r") as f:
            return [c.strip() for c in f.readlines()]

    def init_model(self, config_path: str, weights_path: str):
        net = cv.dnn.readNet(weights_path, config_path)
        net.setPreferableBackend(cv.dnn.DNN_BACKEND_OPENCV)
        net.setPreferableTarget(cv.dnn.DNN_TARGET_CPU)

        model = cv.dnn_DetectionModel(net)
        model.setInputParams(size=(416, 416), scale=1 / 255, swapRB=True)
        return model

    def b64_to_image(self, image_b64: str):
        imgdata = base64.b64decode(str(image_b64))
        img = Image.open(io.BytesIO(imgdata))
        opencv_img = cv.cvtColor(np.array(img), cv.COLOR_BGR2RGB)
        return opencv_img

    def detect(self, frame):
        classes, scores, boxes = self.model.detect(frame, Conf_threshold, NMS_threshold)
        return ODOutput(classes=classes, scores=scores, boxes=boxes)

    def draw_into_frame(self, frame, output: ODOutput):
        for classid, score, box in zip(output.classes, output.scores, output.boxes):
            color = COLORS[int(classid) % len(COLORS)]
            label = "%s : %f" % (self.classes[classid], score)
            cv.rectangle(frame, box, color, 1)
            cv.putText(
                frame,
                label,
                (box[0], box[1] - 10),
                cv.FONT_HERSHEY_COMPLEX,
                0.3,
                color,
                1,
            )
        frame = cv.resize(frame, (240,240), interpolation=cv.INTER_LINEAR)
        return frame

    def run(self, image_b64: str):
        frame = self.b64_to_image(image_b64)
        out = self.detect(frame)
        drawn_frame = self.draw_into_frame(frame, out)
        return drawn_frame
