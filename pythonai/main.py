import pika
import srsly


def callback(channel, method_frame, header_frame, body):
    pass


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
        channel.start_consuming()
    except KeyboardInterrupt:
        channel.stop_consuming()

    connection.close()


if __name__ == "__main__":
    main()
