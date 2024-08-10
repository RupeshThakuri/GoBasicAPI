package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "Book 1", Author: "Author 1", Quantity: 10},
	{ID: "2", Title: "Book 2", Author: "Author 2", Quantity: 5},
	{ID: "3", Title: "Book 3", Author: "Author 3", Quantity: 3},
	{ID: "4", Title: "Book 4", Author: "Author 4", Quantity: 2},
}

// get request function
func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

// post request function
func createBook(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)

	c.IndentedJSON(http.StatusCreated, newBook)
}

// main function to get book by id
func getID(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
	}

	c.IndentedJSON(http.StatusOK, book)
}

// helper functon
func getBookByID(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("book not found")
}

// checkout function like patch (update)
// reduce the quantity by 1
func checkOutBook(c *gin.Context) {
	id, Ok := c.GetQuery("id")
	if !Ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing query parameter id"})
		return
	}

	book, err := getBookByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "book not available"})
	}

	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}

// another patch func that increment quantity by 1
func returnBook(c *gin.Context) {
	id, Ok := c.GetQuery("id")
	if !Ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing query parameter id"})
		return
	}

	book, err := getBookByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

// delete book
func deleteBook(c *gin.Context) {
	id := c.Param("id")
	bookIndex, err := findBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
		return
	}

	books = append(books[:bookIndex], books[bookIndex+1:]...)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "book deleted"})
}

// helper function for delete book
func findBookById(id string) (int, error) {
	for i, b := range books {
		if b.ID == id {
			return i, nil
		}
	}

	return -1, errors.New("book not found")
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", getID)
	router.POST("/books", createBook)
	router.PATCH("/checkout", checkOutBook)
	router.PATCH("/return", returnBook)
	router.DELETE("/books/:id", deleteBook)
	router.Run("localhost:8000")
}
