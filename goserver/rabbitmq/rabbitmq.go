package rabbitmq

import (
	"archome/server/config"
	"archome/server/utils"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", err, msg)
	}
}

func CreateConnection(cfg config.Config) *amqp.Connection {
	r := cfg.RabbitMQ
	connString := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", r.User, r.Pass, r.Host, r.Port, r.VHost)
	conn, err := amqp.Dial(connString)
	failOnError(err, "Failed to create connection")
	return conn
}

func OpenChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open channel")
	return ch
}

func OpenQueue(ch *amqp.Channel) amqp.Queue {
	q, err := ch.QueueDeclare("images", false, false, false, false, nil)
	failOnError(err, "Failed to declare queue")
	return q
}

func PlubishToQueue(ch *amqp.Channel, q amqp.Queue, frame []byte) {

	b64 := utils.EncodeB64(frame)
	err := ch.Publish("", q.Name, false, false, amqp.Publishing{ContentType: "image/jpeg", Body: []byte(b64), Expiration: "600"})
	failOnError(err, "Failed to send message to queue")
}
