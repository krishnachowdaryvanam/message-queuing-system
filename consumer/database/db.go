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

	logrus.Info("Successfully connected to the database")
	return db, nil
}

func ProductExists(db *sql.DB, productID int) error {
	logrus.Info("Checking if product exists for product_id: ", productID)
	productStmt, err := db.Prepare("SELECT COUNT(*) FROM Products WHERE product_id = $1")
	if err != nil {
		logrus.Errorf("Error preparing SQL statement: %v", err)
		return err
	}
	defer productStmt.Close()

	var count int
	err = productStmt.QueryRow(productID).Scan(&count)
	if err != nil {
		logrus.Errorf("Error executing SQL statement: %v", err)
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func GetProductImages(product_id int, db *sql.DB) ([]string, error) {
	err := ProductExists(db, product_id)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Errorf("Product does not exist for product_id: %d", product_id)
			return nil, err
		}
		logrus.Errorf("Error checking if product exists for product_id: %d", product_id)
		return nil, err
	}

	logrus.Info("Getting product images for product_id: ", product_id)
	stmt, err := db.Prepare("SELECT product_images FROM Products WHERE product_id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var product_images string
	err = stmt.QueryRow(product_id).Scan(&product_images)
	if err != nil {
		return nil, err
	}

	images := strings.Split(product_images, ",")
	return images, nil
}

func UpdateProductImages(db *sql.DB, productID int, compressedImagesPaths []string) error {
	err := ProductExists(db, productID)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Errorf("Product does not exist for product_id: %d", productID)
			return err
		}
		logrus.Errorf("Error checking if product exists for product_id: %d", productID)
		return err
	}
	// Update the database
	if len(compressedImagesPaths) == 0 {
		logrus.Error("No images to update")
		return nil
	}
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	compressedImages := strings.Join(compressedImagesPaths, ",")
	query := "UPDATE Products SET compressed_product_images = $1, updated_at = $2 WHERE product_id = $3"
	stmt, err := db.Prepare(query)
	if err != nil {
		logrus.Errorf("error preparing update statement: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(compressedImages, currentTime, productID)
	if err != nil {
		logrus.Errorf("error executing update statement: %v", err)
		return err
	}
	logrus.Infof("Successfully updated product_id: %d", productID)
	return nil
}
