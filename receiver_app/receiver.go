package main

import (
	"log"

	"github.com/fatih/color"
	"github.com/streadway/amqp"
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

	q, err := ch.QueueDeclare(
		"Cola1",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

		log.Print(blue(" ||  Waiting for messages. To exit press CTRL+C"))

	for d := range msgs {
		receivedMessage := green(string(d.Body))
		message := yellow("Received a message: ") + receivedMessage
		log.Print(message)

		// Enviar el mensaje a app3
		err := ch.Publish(
			"",
			"Cola2", // Nombre de la cola para enviar mensajes a app3
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        d.Body,
			},
		)
		failOnError(err, "Failed to publish a message")
	}
}