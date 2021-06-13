package middleware

import (
	"database/sql"
	"encoding/json" // package to encode and decode the json into struct and vice versa
	"fmt"
	"go-postgres/models" // models package where User schema is defined
	"log"
	"net/http" // used to access the request and response object of the api
	"os"       // used to read the environment variable
	"strconv"  // package used to covert string into int type
	"strings"

	"github.com/gorilla/mux" // used to get the params from the route

	"github.com/joho/godotenv" // package used to read the .env file
	_ "github.com/lib/pq"      // postgres golang driver
)

// response format
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// create connection with postgres db
func createConnection() *sql.DB {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Open the connection
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	// return the connection
	return db
}

//Creates a new book object and adds to postgres db
func CreateBook(w http.ResponseWriter, r *http.Request) {

	// set the header to content type x-www-form-urlencoded
	// Allow all origin to handle cors issue
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	//create new book model
	var book models.Book

	//decode the json request to user
	err := json.NewDecoder(r.Body).Decode(&book)

	//check if any errors
	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	//call the insert book function and relay success message
	insertID := insertBook(book)
	message := "Book added successfully"

	//check to see if there was error (Was not sure what/how to handle this the right way so I just checked to see if error id was 400 and if so display message to http)
	if insertID == -1 {
		message = "Rating needs to be in range 1-3"
	}
	//create response object
	res := response{
		ID:      insertID,
		Message: message,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

//Get Book will return book object based on ID
func GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	stringid := strings.ReplaceAll(r.URL.String(), "/api/book/", "")

	// convert the id type from string to int
	id, err := strconv.Atoi(stringid)

	//check if any errors and display message
	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// call the getbookbyID function to get user object and any errors
	book, err := getBookByID(int64(id))

	//check if any errors and display error message
	if err != nil {
		log.Fatalf("Unable to get user. %v", err)
	}

	// send the response
	json.NewEncoder(w).Encode(book)
}

// GetAllBooks will return all the books from database
func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//call get all books method to get all book objects and errors
	books, err := getAllBooks()

	//if there are any errors, display error message
	if err != nil {
		log.Fatalf("Unable to get all user. %v", err)
	}

	// send all the books as response
	json.NewEncoder(w).Encode(books)
}

//function that allows editing/updating of book object information
func UpdateBook(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// get the book id from the request params
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}
	//create new book model
	var book models.Book
	//decode book model data
	err = json.NewDecoder(r.Body).Decode(&book)

	//check if any errors and if so display error message
	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}
	//call update book function that will update book object corresponding to id and new book details
	updatedRows := updateBook(int64(id), book)

	//set message to success message and show how many rows were affected
	msg := fmt.Sprintf("User updated successfully. Total rows/record affected %v", updatedRows)

	//check to see if there was error (Was not sure what/how to handle this the right way so I just checked to see if error id was 400 and if so display message to http)
	if updatedRows == -1 {
		msg = "Rating needs to be in range 1-3"
	}

	//create response object and set error id and message
	res := response{
		ID:      updatedRows,
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

//Delete book will delete book object from database given book id
func DeleteBook(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	stringid := strings.ReplaceAll(r.URL.String(), "/api/deletebook/", "")

	id, err := strconv.Atoi(stringid)

	//check if any errors and return message
	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// call the deletebook function
	deletedRows := deleteBook(int64(id))

	// format the message string
	msg := fmt.Sprintf("User updated successfully. Total rows/record affected %v", deletedRows)

	// format the reponse message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	//send the response
	json.NewEncoder(w).Encode(res)
}

//------------------------- Implementation functions ----------------
//insert book function takes in book model and returns id of book created/inserted
func insertBook(book models.Book) int64 {
	//create connection
	db := createConnection()
	//close the db connection
	defer db.Close()
	//create sql query statement that inserts book into postgres db based on user input data
	sqlStatement := `INSERT INTO book (Title, Author, Publisher, Publish_Date, Rating, Status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING ID`
	//create id variable
	var id int64
	//check to see if rating is within range, if not set return error id of -1 (that way we never actually would return this normally)
	if book.Rating < 1 || book.Rating > 3 {
		return -1
	}
	//query rows based on user input and store into err
	err := db.QueryRow(sqlStatement, book.Title, book.Author, book.Publisher, book.Publish_Date, book.Rating, book.Status).Scan(&id)
	//if there are any errors, return error statement
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	//success message
	fmt.Printf("Inserted a single record %v", id)

	//return the inserted id
	return id
}

// get one book from the DB by its id
func getBookByID(id int64) (models.Book, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create a new book model
	var book models.Book

	// create the select sql query
	sqlStatement := `SELECT * FROM book WHERE id=$1`

	// execute the sql statement
	row := db.QueryRow(sqlStatement, id)

	// unmarshal the row object to book
	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Publisher, &book.Publish_Date, &book.Rating, &book.Status)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return book, nil
	case nil:
		return book, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	return book, err
}

//get every book from database
func getAllBooks() ([]models.Book, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	var books []models.Book

	// create the select sql query
	sqlStatement := `SELECT * FROM book`

	// execute the sql statement
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var book models.Book

		// unmarshal the row object to book
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.Publisher, &book.Publish_Date, &book.Rating, &book.Status)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		//append the book to the books list
		books = append(books, book)

	}

	return books, err
}

// update book from the DB
func updateBook(id int64, book models.Book) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the update sql query
	sqlStatement := `UPDATE book SET Title=$2, Author=$3, Publisher=$4, Publish_Date =$5, Rating = $6, Status = $7 WHERE id=$1`

	//check to see if rating is within correct range and return -1 as error id if out of range
	if book.Rating < 1 || book.Rating > 3 {
		return -1
	}
	// execute the sql statement
	res, err := db.Exec(sqlStatement, id, book.Title, book.Author, book.Publisher, book.Publish_Date, book.Rating, book.Status)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	//check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

// delete book in the DB by id
func deleteBook(id int64) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the delete sql query
	sqlStatement := `DELETE FROM book WHERE id=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}
