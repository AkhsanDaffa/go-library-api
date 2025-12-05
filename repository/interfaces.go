package repository

import "go-library-api/models"

type UserRepositoryInterface interface {
	CreateUser(username, email, hashedPassword string) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
}
