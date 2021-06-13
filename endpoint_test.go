package main

import (
	"bytes" // package to encode and decode the json into struct and vice versa
	"fmt"

	// models package where User schema is defined
	"go-postgres/middleware"
	"net/http"          // used to access the request and response object of the api
	"net/http/httptest" // used to read the environment variable
	"strings"
	"testing"

	// package used to covert string into int type
	// used to get the params from the route

	// package used to read the .env file
	_ "github.com/lib/pq"
)

func TestGetBooks(t *testing.T) {

	req, err := http.NewRequest("GET", "/api/book", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(middleware.GetAllBooks)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	fmt.Print(req.URL)
	// Check the response body is what we expect.
	expected := `[{"ID":8,"Title":"Harry Potter and the Chamber of Secrets","Author":"J.K. Rowling","Publisher":"Bloomsbury","Publish_Date":"1998-07-02T00:00:00Z","Rating":3,"Status":false},{"ID":9,"Title":"Harry Potter and the Chamber of Secrets","Author":"J.K. Rowling","Publisher":"Bloomsbury","Publish_Date":"1998-07-02T00:00:00Z","Rating":3,"Status":false},{"ID":10,"Title":"Harry Potter and the Chamber of Secrets","Author":"J.K. Rowling","Publisher":"Bloomsbury","Publish_Date":"1998-07-02T00:00:00Z","Rating":2,"Status":false}]`
	expected = strings.TrimRight(expected, "\r\n")
	data := strings.TrimRight(rr.Body.String(), "\r\n")
	if data != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			data, expected)
	}
}
func TestGetBook(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/book/20", nil)
	//req, err := http.NewRequest("GET", "/api/book/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(req.URL, " ")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(middleware.GetBook)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"ID":20,"Title":"Harry Potter and the Chamber of Secrets","Author":"J.K. Rowling","Publisher":"Bloomsbury","Publish_Date":"1998-07-02T00:00:00Z","Rating":3,"Status":false}`
	expected = strings.TrimRight(expected, "\r\n")
	data := strings.TrimRight(rr.Body.String(), "\r\n")
	if data != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
func TestGetBookFailed(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/book/1", nil)
	//req, err := http.NewRequest("GET", "/api/book/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(req.URL, " ")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(middleware.GetBook)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"ID":0,"Title":"","Author":"","Publisher":"","Publish_Date":"","Rating":0,"Status":false}`
	expected = strings.TrimRight(expected, "\r\n")
	data := strings.TrimRight(rr.Body.String(), "\r\n")
	if data != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
func TestCreateBook(t *testing.T) {

	var jsonStr = []byte(`{"Title":"Harry Potter and the Chamber of Secrets","Author":"J.K. Rowling","Publisher":"Bloomsbury","Publish_Date":"1998-07-02T00:00:00Z","Rating":3,"Status":false}`)

	req, err := http.NewRequest("POST", "/api/newbook", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(middleware.CreateBook)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `"message":"Book added successfully"`
	expected = strings.TrimRight(expected, "\r\n")
	data := strings.TrimRight(rr.Body.String(), "\r\n")

	if !strings.Contains(data, expected) {
		t.Errorf("Handler response %v did not contain %v",
			data, expected)
	}

}

func TestDeleteBook(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/api/deletebook/10", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(middleware.DeleteBook)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"id":10,"message":"User updated successfully. Total rows/record affected 1"}`
	expected = strings.TrimRight(expected, "\r\n")
	data := strings.TrimRight(rr.Body.String(), "\r\n")
	if data != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestDeleteBookFailed(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/api/deletebook/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(middleware.DeleteBook)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"id":1,"message":"User updated successfully. Total rows/record affected 0"}`
	expected = strings.TrimRight(expected, "\r\n")
	data := strings.TrimRight(rr.Body.String(), "\r\n")
	if data != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
