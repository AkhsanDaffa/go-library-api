package repository

import (
	"database/sql"
	"go-library-api/models"
)

type CheckoutRepository struct {
	DB *sql.DB
}

func NewCheckoutRepository(db *sql.DB) *CheckoutRepository {
	return &CheckoutRepository{DB: db}
}

func (r *CheckoutRepository) CreateCheckout(UserID int, BookID int) (models.Checkout, error) {
	var checkout models.Checkout

	query := `
	INSERT INTO checkouts (user_id, book_id)
	VALUES ($1, $2)
	RETURNING id, checkout_date, return_date, user_id, book_id
	`

	err := r.DB.QueryRow(query, UserID, BookID).Scan(
		&checkout.ID,
		&checkout.CheckoutDate,
		&checkout.ReturnDate,
		&checkout.UserID,
		&checkout.BookID,
	)

	if err != nil {
		return models.Checkout{}, err
	}

	return checkout, nil
}

func (r *CheckoutRepository) GetAllCheckouts() ([]models.Checkout, error) {
	var checkouts []models.Checkout

	query := `
	SELECT id, checkout_date, return_date, user_id, book_id FROM checkouts
	`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var checkout models.Checkout

		if err := rows.Scan(&checkout.ID, &checkout.CheckoutDate, &checkout.ReturnDate, &checkout.UserID, &checkout.BookID); err != nil {
			return nil, err
		}

		checkouts = append(checkouts, checkout)
	}
	return checkouts, nil
}

func (r *CheckoutRepository) GetCheckoutByID(id int) (models.Checkout, error) {
	var checkout models.Checkout

	query := `
	SELECT id, checkout_date, return_date, user_id, book_id FROM checkouts WHERE id = $1
	`

	err := r.DB.QueryRow(query, id).Scan(
		&checkout.ID,
		&checkout.CheckoutDate,
		&checkout.ReturnDate,
		&checkout.UserID,
		&checkout.BookID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Checkout{}, nil
		}
		return models.Checkout{}, err
	}
	return checkout, nil
}
