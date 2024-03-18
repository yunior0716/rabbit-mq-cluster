package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
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

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/rabbitmq_go")
	failOnError(err, "Failed to connect to MySQL")
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS messages (id INT AUTO_INCREMENT PRIMARY KEY, content TEXT, received_at DATETIME)")
	failOnError(err, "Failed to create table")

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

	log.Println(" [*] Waiting for messages. To exit press CTRL+C")

	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)

		// Obtener la hora actual
		receivedAt := time.Now().Format("2006-01-02 15:04:05")

		// Guardar el mensaje en MySQL junto con la hora actual
		_, err := db.Exec("INSERT INTO messages (content, received_at) VALUES (?, ?)", string(d.Body), receivedAt)
		failOnError(err, "Failed to insert message into MySQL")

		log.Println("Message saved in MySQL")
	}
}