package getCSVRoutes

import (
	"database/sql"
	"encoding/csv"
	"io"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

var DB *sql.DB

type Book struct {
	ID               int
	Title            string
	Author           string
	Publication_year int
}

// GetCSVRoutes is an exported function that sets up and returns Fiber app with CSV routes
func SetupCSVRoutes(app *fiber.App, DataBase *sql.DB) {

	DB = DataBase

	csvGroup := app.Group("/csv")

	csvGroup.Put("/updateFromCSV", uploadCSVHandler)

	csvGroup.Get("/read", fetchDataHandler)

}

func fetchDataHandler(c *fiber.Ctx) error {
	rows, err := DB.Query("SELECT id, title, author, publication_year FROM books")
	if err != nil {
		log.Println("Error executing query:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	defer rows.Close()

	var data []Book

	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Publication_year); err != nil {
			log.Println("Error scanning row:", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		data = append(data, book)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(data)
}
func uploadCSVHandler(c *fiber.Ctx) error {

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("File upload error")
	}

	fileContent, err := file.Open()
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Error opening file")
	}
	defer fileContent.Close()

	reader := csv.NewReader(fileContent)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Error reading CSV")
		}

		for _, field := range record {
			if field == "" {
				return c.Status(http.StatusBadRequest).SendString("Enter Some Data")
			}
		}

		_, err = DB.Exec("INSERT INTO books (ID, title, author, publication_year) VALUES ($1, $2, $3, $4)",
			record[0], record[1], record[2], record[3])

		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Error inserting data into the database")
		}
	}

	return c.SendString("CSV data uploaded and inserted into the database")
}
