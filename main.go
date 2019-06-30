package main

import (
	"books-list/controllers"
	"books-list/driver"
	"books-list/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/subosito/gotenv"

	"github.com/gorilla/mux"
)

type Book struct {
	ID     string `json:id`
	Title  string `json:title`
	Author string `json:author`
	Year   string `json:year`
}

var books []models.Book

var db *sql.DB

func init() {
	gotenv.Load()
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func main() {

	db = driver.ConnectDB()

	controller := controllers.Controller{}
	router := mux.NewRouter()

	router.HandleFunc("/books", controller.GetBooks(db)).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/books", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", removeBook).Methods("DELETE")

	fmt.Println("Server is running at port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func getBook(w http.ResponseWriter, r *http.Request) {
	log.Println("Get book by ID is called")
	var book models.Book
	params := mux.Vars(r)
	fmt.Println(params)

	rows := db.QueryRow("Select * from books where id=$1", params["id"])

	err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	logFatal(err)

	json.NewEncoder(w).Encode(book)
}
func addBook(w http.ResponseWriter, r *http.Request) {
	log.Println("add  book is called")

	var book models.Book
	var bookID int

	json.NewDecoder(r.Body).Decode(&book)

	err := db.QueryRow("insert into books(title,author,year) values($1,$2,$3) RETURNING id;",
		book.Title, book.Author, book.Year).Scan(&book.ID)

	logFatal(err)

	json.NewEncoder(w).Encode(bookID)

}
func updateBook(w http.ResponseWriter, r *http.Request) {
	log.Println("Update book is called")
	var book models.Book
	json.NewDecoder(r.Body).Decode(&book)

	result, err := db.Exec("update books set title=$1,author=$2,year=$3 where id=$4 RETURNING id",
		&book.Title, &book.Author, &book.Year, &book.ID)
	logFatal(err)
	rowsUpdated, err := result.RowsAffected()
	logFatal(err)
	json.NewEncoder(w).Encode(rowsUpdated)
}
func removeBook(w http.ResponseWriter, r *http.Request) {
	log.Println("Remove book is called")
	params := mux.Vars(r)

	result, err := db.Exec("Delete from books where id = $1", params["id"])
	logFatal(err)
	rowsDeleted, err := result.RowsAffected()
	logFatal(err)

	json.NewEncoder(w).Encode(rowsDeleted)
}
