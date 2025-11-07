package repository

import (
	"database/sql"
	"go-library-api/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(username string, email string) (models.User, error) {
	var user models.User

	query := `
	INSERT INTO users (username, email)
	VALUES ($1, $2)
	RETURNING id, username, email, created_at
	`

	err := r.DB.QueryRow(query, username, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User

	query := `
	SELECT id, username, email, created_at FROM users
	`
	rows, err := r.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user models.User

		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt); err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) GetUserByID(id int) (models.User, error) {
	var user models.User

	query := `
	SELECT id, username, email, created_at FROM users WHERE id = $1
	`
	err := r.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, nil
		}
		return models.User{}, err
	}
	return user, nil
}
