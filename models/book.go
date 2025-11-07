package models

type Book struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	IsAvailable bool   `json:"is_available"`
}
