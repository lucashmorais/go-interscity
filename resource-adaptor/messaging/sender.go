package messaging

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) bool {
	if err != nil {
		// log.Fatalf("%s: %s", msg, err)
		log.Printf("%s: %s", msg, err)
		return true
	}
	return false
}

func coreSend(ch *amqp.Channel, q amqp.Queue, body string) {
	for range time.Tick(time.Second * 1) {
		err := ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		if failOnError(err, "Failed to publish a message") {
			return
		}
	}
}

func ConnectAndSend() {
	conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Initializign connection
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Creating queue
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := "Hello World!"

	// var wg sync.WaitGroup
	go coreSend(ch, q, body)
	time.Sleep(time.Second * 5)
}
