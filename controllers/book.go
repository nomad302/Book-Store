package controllers

import (
	"books-list/models"
	bookRepository "books-list/repository/book"
	"books-list/utils"
	"database/sql"
	"log"
	"net/http"
)

type Controller struct{}

var books []models.Book

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// returning anonymous function
func (c Controller) GetBooks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Get all books is called")
		var book models.Book
		books = []models.Book{}
		var error models.Error

		bookRepo := bookRepository.BookRepository{}

		books, err := bookRepo.GetBooks(db, book, books)

		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		utils.SendSuccess(w, books)
	}
}
