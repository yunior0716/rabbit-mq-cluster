package main

import (
	"log"

	"github.com/streadway/amqp"

	"github.com/fatih/color"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {

	//red := color.New(color.FgRed).SprintfFunc()
	green := color.New(color.FgGreen).SprintfFunc()
	blue := color.New(color.FgBlue).SprintfFunc()
	yellow := color.New(color.FgYellow).SprintfFunc()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q1, err := ch.QueueDeclare(
		"Cola1",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare queue 1")

	body := "Mensaje de prueba 1 : Delio Diaz"
	sentMessage := green(body)
	message := yellow("Sent: ") + sentMessage
	err = ch.Publish(
		"",
		q1.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	failOnError(err, "Failed to publish a message")
	log.Print(message)

	q3, err := ch.QueueDeclare(
		"Cola3",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare queue")

	msgs, err := ch.Consume(
		q3.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer for queue 3")

	log.Print(blue(" ||  Waiting for messages. To exit press CTRL+C"))

	for d := range msgs {
		receivedMessage := green(string(d.Body))
		log.Printf("Received a message: %s", receivedMessage)
	}
}