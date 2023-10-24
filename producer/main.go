package main

import (
	"log"
	"messagequeuesystem/producer/database"
	"messagequeuesystem/producer/handlers"
	"messagequeuesystem/producer/message"

	"github.com/gin-gonic/gin"
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

	queue := "products"

	conn, err := message.NewRMQ()
	if err != nil {
		logrus.Errorf("Failed to connect to RabbitMQ: %v", err)
		return
	}
	defer conn.Close()

	ch, err := message.NewChannel(conn)
	if err != nil {
		logrus.Errorf("Failed to open an RMQ channel: %v", err)
		return
	}
	defer ch.Close()

	//Gin router
	router := gin.Default()

	//route to receive the product data
	router.POST("/products", handlers.SaveProduct(db, ch, queue))

	// Start the Gin server
	if err := router.Run(":8080"); err != nil {
		logrus.Fatalf("Error in starting the server: %v", err)
	}
}
