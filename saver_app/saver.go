package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/fatih/color"
	_ "github.com/go-sql-driver/mysql"
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

	// Conexión a RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Canal de comunicación con RabbitMQ
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declaración de la cola para recibir mensajes de app2
	q, err := ch.QueueDeclare(
		"Cola2",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	// Conexión a la base de datos MySQL
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/rabbitmq_go")
	failOnError(err, "Failed to connect to MySQL")
	defer db.Close()

	// Crear tabla si no existe
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS messages (id INT AUTO_INCREMENT PRIMARY KEY, content TEXT, received_at DATETIME)")
	failOnError(err, "Failed to create table")

	// Consumir mensajes de la cola
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

	// Procesar mensajes recibidos
	for d := range msgs {
		receivedMessage := green(string(d.Body))
		message := yellow("Received a message: ") + receivedMessage
		log.Print(message)

		// Obtener la hora actual
		receivedAt := time.Now().Format("2006-01-02 15:04:05")

		// Guardar el mensaje en la base de datos junto con la hora actual
		_, err := db.Exec("INSERT INTO messages (content, received_at) VALUES (?, ?)", string(d.Body), receivedAt)
		failOnError(err, "Failed to insert message into MySQL")

		log.Println(blue("Message saved successfully"))

		// Enviar un mensaje a app1 para informar que el mensaje fue guardado correctamente
		err = ch.Publish(
			"",
			"Cola3", // Nombre de la cola para enviar mensajes a app1
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte("Message saved successfully"),
			},
		)
		failOnError(err, "Failed to publish a message")
	}
}