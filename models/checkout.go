package models

import "time"

type Checkout struct {
	ID           int       `json:"id"`
	CheckoutDate time.Time `json:"checkout_date"`
	ReturnDate   time.Time `json:"return_date"`
	UserID       int       `json:"user_id"`
	BookID       int       `json:"book_id"`
}
