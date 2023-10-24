package main

import (
	"log"
	"messagequeuesystem/consumer/database"
	"messagequeuesystem/consumer/message"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize the database connection.
	db, err := database.InitializeDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize the database: %s", err)
	}
	defer db.Close()

	queue := "products" // Set the queue name directly in your code

	conn, err := message.NewRMQ()
	if err != nil {
		logrus.Errorf("Failed to connect to RabbitMQ: %v", err)
		return
	}
	defer conn.Close()
	ch, err := message.NewChannel(conn)
	if err != nil {
		logrus.Errorf("Failed to open a rmq channel: %v", err)
		return
	}
	defer ch.Close()
	image_quality := 60
	message.Consumer(ch, queue, db, image_quality)
}
