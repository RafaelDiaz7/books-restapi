package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

//Book Struct (Model)
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
	Pages int32	`json:"pages"`
}

type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

//Init books var as slice book struct
var books []Book

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:3000")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:3000")
	params := mux.Vars(r) //get params
	//Loop through books and find with id
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

func createBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:3000")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(10000000)) // Mock iD - not safe in production
	books = append(books, book)
	json.NewEncoder(w).Encode(book)

	/* what should happen is it add the book
	   and there's no database is just gonna
	   be adding it basically to memory and
	   then it will output a response of that
	   new book
	*/
}

func updateBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:3000")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = strconv.Itoa(rand.Intn(10000000)) // Mock iD - not safe in production
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}

	}
	json.NewEncoder(w).Encode(books)
}

func deleteBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:3000")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
		}
		break
	}
	json.NewEncoder(w).Encode(books)
}

func main() {
	//Init port
	port := os.Getenv("PORT")

	//Init router
	r := mux.NewRouter()

	//MOck data - @todo - implement db
	books = append(books, Book{ID: "1", Isbn: "448743", Title: "Book One", Author: &Author{Firstname: "John", Lastname: "Doe"},Pages:560})
	books = append(books, Book{ID: "2", Isbn: "848893", Title: "Book Two", Author: &Author{Firstname: "Tyler", Lastname: "Mcay"},Pages:978})

	//Route Handlers / Endpoints
	r.HandleFunc("/api/books", getBooks).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/books", createBooks).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/books/{id}", updateBooks).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/books", deleteBooks).Methods("DELETE", "OPTIONS")

	// Using cors library to enable access control allow headers
	handler := cors.Default().Handler(r)

	log.Fatal(http.ListenAndServe(":"+port, handler))
}
