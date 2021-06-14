package main

import (
	"bytes"
	"go-postgres/middleware"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/lib/pq"
)

func TestSetUp(t *testing.T) {
	//Deletes all entries in book table and resets primary key sequence
	middleware.PrepForTesting()
}
func TestCreateBook(t *testing.T) {
	//tests adding a new book
	var jsonStr = []byte(`{"Title":"Crime and Punishment","Author":"Fyodor Dostoyevsky","Publisher":"The Russian Messenger","Publish_Date":"1886-02-15","Rating":2.8,"Status":false}`)
	//create new POST request using api route for adding a new book
	req, err := http.NewRequest("POST", "/api/newbook", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	//set required headers
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(middleware.CreateBook)
	handler.ServeHTTP(rr, req)
	//check to see if connected properly
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	//should return success message
	expected := `"message":"Book added successfully"`
	expected = strings.TrimRight(expected, "\r\n")
	data := strings.TrimRight(rr.Body.String(), "\r\n")
	//check to see if correct message was returned successfully
	if !strings.Contains(data, expected) {
		t.Errorf("Handler response %v did not contain %v",
			data, expected)
	}
	//Test adding another book
	jsonStr = []byte(`{"Title":"Harry Potter and the Chamber of Secrets","Author":"J.K. Rowling","Publisher":"Bloomsbury","Publish_Date":"1998-07-02T00:00:00Z","Rating":3,"Status":false}`)

	req, err = http.NewRequest("POST", "/api/newbook", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(middleware.CreateBook)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected = `"message":"Book added successfully"`
	expected = strings.TrimRight(expected, "\r\n")
	data = strings.TrimRight(rr.Body.String(), "\r\n")

	if !strings.Contains(data, expected) {
		t.Errorf("Handler response %v did not contain %v",
			data, expected)
	}

}
func TestCreateBookFailed(t *testing.T) {
	//First test checks if user tries to update book with rating above the accepted range
	var jsonStr = []byte(`{"Title":"Harry Potter and the Chamber of Secrets","Author":"J.K. Rowling","Publisher":"Bloomsbury","Publish_Date":"1998-07-02T00:00:00Z","Rating":5,"Status":true}`)

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
	//should show message correct rating range to user
	expected := `{"id":-1,"message":"Rating needs to be in range 1-3"}`
	expected = strings.TrimRight(expected, "\r\n")
	data := strings.TrimRight(rr.Body.String(), "\r\n")
	if data != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
	//Third test checks if user tries to update book with rating bellow the accepted range
	jsonStr = []byte(`{"Title":"Harry Potter and the Chamber of Secrets","Author":"J.K. Rowling","Publisher":"Bloomsbury","Publish_Date":"1998-07-02T00:00:00Z","Rating":0,"Status":true}`)

	req, err = http.NewRequest("POST", "/api/newbook", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(middleware.CreateBook)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	//should show message correct rating range to user
	expected = `{"id":-1,"message":"Rating needs to be in range 1-3"}`
	expected = strings.TrimRight(expected, "\r\n")
	data = strings.TrimRight(rr.Body.String(), "\r\n")
	if data != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

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

	// Check the response body is what we expect.
	expected := `[{"ID":1,"Title":"Crime and Punishment","Author":"Fyodor Dostoyevsky","Publisher":"The Russian Messenger","Publish_Date":"1886-02-15T00:00:00Z","Rating":2.8,"Status":false},{"ID":2,"Title":"Harry Potter and the Chamber of Secrets","Author":"J.K. Rowling","Publisher":"Bloomsbury","Publish_Date":"1998-07-02T00:00:00Z","Rating":3,"Status":false}]`
	expected = strings.TrimRight(expected, "\r\n")
	data := strings.TrimRight(rr.Body.String(), "\r\n")
	if data != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			data, expected)
	}
}
func TestGetBook(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/book/2", nil)
	//req, err := http.NewRequest("GET", "/api/book/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(middleware.GetBook)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"ID":2,"Title":"Harry Potter and the Chamber of Secrets","Author":"J.K. Rowling","Publisher":"Bloomsbury","Publish_Date":"1998-07-02T00:00:00Z","Rating":3,"Status":false}`
	expected = strings.TrimRight(expected, "\r\n")
	data := strings.TrimRight(rr.Body.String(), "\r\n")
	if data != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
func TestGetBookFailed(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/book/3", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(middleware.GetBook)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"ID":0,"Title":"","Author":"","Publisher":"","Publish_Date":"","Rating":0,"Status":false}`
	expected = strings.TrimRight(expected, "\r\n")
	data := strings.TrimRight(rr.Body.String(), "\r\n")
	if data != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
func TestEditEntry(t *testing.T) {

	var jsonStr = []byte(`{"ID":2,"Title":"Harry Potter and the Chamber of Secrets","Author":"J.K. Rowling","Publisher":"Bloomsbury","Publish_Date":"1998-07-02T00:00:00Z","Rating":2.65,"Status":true}`)

	req, err := http.NewRequest("PUT", "/api/book/2", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(middleware.UpdateBook)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"id":2,"message":"User updated successfully. Total rows/record affected 1 "}`
	expected = strings.TrimRight(expected, "\r\n")
	data := strings.TrimRight(rr.Body.String(), "\r\n")
	if data != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
func TestEditEntryFailed(t *testing.T) {
	//first test for trying to edit book that doesnt exist
	var jsonStr = []byte(`{"ID":2,"Title":"Harry Potter and the Chamber of Secrets","Author":"J.K. Rowling","Publisher":"Bloomsbury","Publish_Date":"1998-07-02T00:00:00Z","Rating":2.65,"Status":true}`)

	req, err := http.NewRequest("PUT", "/api/book/3", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(middleware.UpdateBook)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	//should show that 0 rows/records were affected because book with id 3 does not exist
	expected := `{"id":3,"message":"User updated successfully. Total rows/record affected 0 "}`
	expected = strings.TrimRight(expected, "\r\n")
	data := strings.TrimRight(rr.Body.String(), "\r\n")
	if data != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	//Second test checks if user tries to update book with rating above the accepted range
	jsonStr = []byte(`{"ID":2,"Title":"Harry Potter and the Chamber of Secrets","Author":"J.K. Rowling","Publisher":"Bloomsbury","Publish_Date":"1998-07-02T00:00:00Z","Rating":5,"Status":true}`)

	req, err = http.NewRequest("PUT", "/api/book/2", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(middleware.UpdateBook)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	//should show message correct rating range to user
	expected = `{"id":2,"message":"Rating needs to be in range 1-3"}`
	expected = strings.TrimRight(expected, "\r\n")
	data = strings.TrimRight(rr.Body.String(), "\r\n")
	if data != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
	//Third test checks if user tries to update book with rating bellow the accepted range
	jsonStr = []byte(`{"ID":2,"Title":"Harry Potter and the Chamber of Secrets","Author":"J.K. Rowling","Publisher":"Bloomsbury","Publish_Date":"1998-07-02T00:00:00Z","Rating":0,"Status":true}`)

	req, err = http.NewRequest("PUT", "/api/book/2", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(middleware.UpdateBook)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	//should show message correct rating range to user
	expected = `{"id":2,"message":"Rating needs to be in range 1-3"}`
	expected = strings.TrimRight(expected, "\r\n")
	data = strings.TrimRight(rr.Body.String(), "\r\n")
	if data != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
func TestDeleteBook(t *testing.T) {
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
	expected := `{"id":1,"message":"User updated successfully. Total rows/record affected 1"}`
	expected = strings.TrimRight(expected, "\r\n")
	data := strings.TrimRight(rr.Body.String(), "\r\n")
	if data != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestDeleteBookFailed(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/api/deletebook/3", nil)
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
	expected := `{"id":3,"message":"User updated successfully. Total rows/record affected 0"}`
	expected = strings.TrimRight(expected, "\r\n")
	data := strings.TrimRight(rr.Body.String(), "\r\n")
	if data != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
