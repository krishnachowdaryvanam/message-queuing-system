package database

import (
	"database/sql"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/sirupsen/logrus"
)

// InitializeDatabase opens a database connection to PostgreSQL.
func InitializeDatabase() (*sql.DB, error) {
	connStr := "user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Automatically create the necessary tables if they don't exist
	err = createTablesIfNotExist(db)
	if err != nil {
		return nil, err
	}
	logrus.Info("Successfully connected to the database")
	return db, nil
}

func createTablesIfNotExist(db *sql.DB) error {
	// Define SQL statements to create tables if they don't exist
	createUsersTableSQL := `
	CREATE TABLE IF NOT EXISTS Users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255),
		mobile VARCHAR(20),
		latitude DOUBLE PRECISION,
		longitude DOUBLE PRECISION,
		created_at TIMESTAMP,
		updated_at TIMESTAMP
	);
	`

	createProductsTableSQL := `
	CREATE TABLE IF NOT EXISTS Products (
		product_id SERIAL PRIMARY KEY,
		product_name VARCHAR(255),
		product_description TEXT,
		product_images TEXT,
		product_price DECIMAL(10, 2),
		compressed_product_images TEXT,
		created_at TIMESTAMP,
		updated_at TIMESTAMP
	);
	`

	// Execute the SQL statements to create tables
	_, err := db.Exec(createUsersTableSQL)
	if err != nil {
		return err
	}

	_, err = db.Exec(createProductsTableSQL)
	if err != nil {
		return err
	}
	logrus.Info("Successfully created tables")
	return nil
}

func UserExists(db *sql.DB, userID int) error {
	logrus.Info("Checking if user exists for user_id: ", userID)
	userStmt, err := db.Prepare("SELECT COUNT(*) FROM Users WHERE id = $1")
	if err != nil {
		logrus.Errorf("Error preparing SQL statement: %v", err)
		return err
	}
	defer userStmt.Close()

	var count int
	err = userStmt.QueryRow(userID).Scan(&count)
	if err != nil {
		logrus.Errorf("Error executing SQL statement: %v", err)
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func InsertProduct(db *sql.DB, ProductName string, ProductDescription string, ProductPrice float64, productImages []string) (int64, error) {
	// Join the product images into a comma-separated string
	productImagesStr := strings.Join(productImages, ",")

	// Insert the product into the database
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	stmt, err := db.Prepare("INSERT INTO Products (product_name, product_description, product_images, product_price, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING product_id")
	if err != nil {
		logrus.Errorf("Error preparing SQL statement: %v", err)
		return 0, err
	}
	defer stmt.Close()

	var productID int64
	err = stmt.QueryRow(ProductName, ProductDescription, productImagesStr, ProductPrice, currentTime).Scan(&productID)
	if err != nil {
		logrus.Errorf("Error executing SQL statement: %v", err)
		return 0, err
	}

	logrus.Infof("Successfully inserted product into the database with ID: %d", productID)
	return productID, nil
}
