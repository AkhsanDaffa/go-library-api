package repository

import (
	"database/sql"
	"go-library-api/models"
)

type BookRepository struct {
	DB *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{DB: db}
}

func (r *BookRepository) CreateBook(title string, author string, is_available bool) (models.Book, error) {
	var book models.Book

	query := `
	INSERT INTO books (title, author, is_available)
	VALUES ($1, $2, $3)
	RETURNING id, title, author, is_available
	`

	err := r.DB.QueryRow(query, title, author, is_available).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.IsAvailable,
	)

	if err != nil {
		return models.Book{}, err
	}

	return book, nil
}

func (r *BookRepository) GetAllBooks() ([]models.Book, error) {
	var books []models.Book

	query := `
	SELECT id, title, author, is_available FROM books
	`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var book models.Book

		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.IsAvailable); err != nil {
			return nil, err
		}

		books = append(books, book)
	}
	return books, nil
}

func (r *BookRepository) GetBookByID(id int) (models.Book, error) {
	var book models.Book

	query := `
	SELECT id, title, author, is_available FROM books WHERE id = $1
	`

	err := r.DB.QueryRow(query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.IsAvailable,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Book{}, nil
		}

		return models.Book{}, err
	}
	return book, nil
}

func (r *BookRepository) UpdateBookAvailability(bookID int, isAvailable bool) error {
	query := `
	UPDATE books SET is_available = $1 WHERE id = $2
	`

	_, err := r.DB.Exec(query, isAvailable, bookID)
	return err
}
