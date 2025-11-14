package models

import (
	"database/sql"
	"time"
)

type Checkout struct {
	ID           int          `json:"id"`
	UserID       int          `json:"user_id"`
	BookID       int          `json:"book_id"`
	CheckoutDate time.Time    `json:"checkout_date"`
	ReturnDate   sql.NullTime `json:"return_date"`
}
