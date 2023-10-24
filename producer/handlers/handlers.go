package handlers

import (
	"database/sql"
	"log"
	"messagequeuesystem/producer/database"
	"messagequeuesystem/producer/message"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

type Product struct {
	UserID             int      `json:"user_id"`
	ProductName        string   `json:"product_name"`
	ProductDescription string   `json:"product_description"`
	ProductImages      []string `json:"product_images"`
	ProductPrice       float64  `json:"product_price"`
}

func SaveProduct(db *sql.DB, ch *amqp.Channel, queue string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the request body into a Product struct
		var product Product
		if err := c.ShouldBindJSON(&product); err != nil {
			log.Printf("Error in parsing the request body: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		err := database.UserExists(db, product.UserID)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("User not found: %v", err)
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			} else {
				log.Printf("Error in checking if user exists: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}
		}

		productID, err := database.InsertProduct(db, product.ProductName, product.ProductDescription, product.ProductPrice, product.ProductImages)
		if err != nil {
			log.Printf("Error in inserting product: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		err = message.Producer(productID, ch, queue)
		if err != nil {
			log.Printf("Error in sending message to queue: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		// Return a success message
		c.String(http.StatusOK, "Product saved successfully")
	}
}
