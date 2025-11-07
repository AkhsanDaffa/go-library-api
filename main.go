package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"

	"go-library-api/config"
	"go-library-api/handler"
	"go-library-api/repository"
)

func main() {
	log.Println("Memulai Aplikasi Perpustakaan")

	connStr := config.GetDBConnectionString()

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Failed Open DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed Connect DB: %v", err)
	}

	log.Println("Connect DB Successful")

	userRepo := repository.NewUserRepository(db)
	bookRepo := repository.NewBookRepository(db)
	checkoutRepo := repository.NewCheckoutRepository(db)

	userHandler := handler.NewUserHandler(userRepo)
	bookHandler := handler.NewBookHandler(bookRepo)
	checkoutHandler := handler.NewCheckoutHandler(checkoutRepo, bookRepo)

	r := gin.Default()

	userRoutes := r.Group("/users")
	{
		// POST /users
		userRoutes.POST("", userHandler.CreateUser)

		// GET /users
		userRoutes.GET("", userHandler.GetAllUsers)

		// GET /users/:id
		userRoutes.GET("/:id", userHandler.GetUserByID)
	}

	bookRoutes := r.Group("/books")
	{
		// POST /books
		bookRoutes.POST("", bookHandler.CreateBook)
		// GET /books
		bookRoutes.GET("", bookHandler.GetAllBooks)

		// GET /books/:id
		bookRoutes.GET("/:id", bookHandler.GetBookByID)
	}

	checkoutRoutes := r.Group("/checkouts")
	{
		// POST /checkouts
		checkoutRoutes.POST("", checkoutHandler.CheckoutBook)
	}

	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
