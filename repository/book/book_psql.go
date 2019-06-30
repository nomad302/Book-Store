package bookRepository

import (
	"books-list/models"
	"database/sql"
	"log"
)

type BookRepository struct{}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (b BookRepository) GetBooks(db *sql.DB, book models.Book, books []models.Book) ([]models.Book, error) {
	rows, err := db.Query("Select * from books")

	if err != nil {
		return []models.Book{}, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		logFatal(err)
		books = append(books, book)
	}

	if err != nil {
		return []models.Book{}, err
	}

	return books, nil
}
