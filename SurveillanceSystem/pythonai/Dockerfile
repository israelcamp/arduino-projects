FROM python:3.13
   
RUN pip install uv
RUN apt-get update && apt-get install ffmpeg libsm6 libxext6  -y

WORKDIR /app

COPY ./pyproject.toml ./uv.lock ./
RUN uv sync 

COPY ./yolo main.py ./

