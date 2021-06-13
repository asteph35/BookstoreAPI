package models

// User schema of the user table
type Book struct {
	ID           int64   `json:"ID"`
	Title        string  `json:"Title"`
	Author       string  `json:"Author"`
	Publisher    string  `json:"Publisher"`
	Publish_Date string  `json:"Publish_Date"`
	Rating       float64 `json:"Rating"`
	Status       bool    `json:"Status"`
}
