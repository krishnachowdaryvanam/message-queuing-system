package database

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func TestInsertProduct(t *testing.T) {
	// Set up a test database
	testDB, err := sql.Open("postgres", ":memory:")
	if err != nil {
		t.Fatalf("Error opening test database: %v", err)
	}
	defer testDB.Close()

	// Create the Products table
	_, err = testDB.Exec(`
        CREATE TABLE Products (
            id SERIAL PRIMARY KEY,
            product_name VARCHAR(255) NOT NULL,
            product_description TEXT NOT NULL,
            product_images TEXT[] NOT NULL,
            product_price DECIMAL(10,2) NOT NULL,
            created_at TIMESTAMP NOT NULL
        )
    `)
	if err != nil {
		t.Fatalf("Error creating Products table: %v", err)
	}

	// Call the function to insert a product into the test database
	productImages := []string{"image1.jpg", "image2.jpg"}
	productID, err := InsertProduct(testDB, "Test Product", "A test product", 9.99, productImages)
	if err != nil {
		t.Fatalf("Error inserting product: %v", err)
	}

	// Check that the product was inserted with the correct values
	var productName, productDescription, productImagesStr string
	var productPrice float64
	var createdAt time.Time
	err = testDB.QueryRow("SELECT * FROM Products WHERE id = $1", productID).Scan(&productID, &productName, &productDescription, &productImagesStr, &productPrice, &createdAt)
	if err != nil {
		t.Fatalf("Error querying product: %v", err)
	}
	if productName != "Test Product" {
		t.Errorf("Expected product_name to be 'Test Product', but got '%s'", productName)
	}
	if productDescription != "A test product" {
		t.Errorf("Expected product_description to be 'A test product', but got '%s'", productDescription)
	}
	if productImagesStr != "image1.jpg,image2.jpg" {
		t.Errorf("Expected product_images to be 'image1.jpg,image2.jpg', but got '%s'", productImagesStr)
	}
	if productPrice != 9.99 {
		t.Errorf("Expected product_price to be 9.99, but got %f", productPrice)
	}
}
