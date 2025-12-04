package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"

	"go-library-api/config"
	"go-library-api/handler"
	"go-library-api/middleware"
	"go-library-api/repository"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "go-library-api/docs"
)

// @title Library API
// @version 1.0
// @description Aplikasi Manajemen Perpustakaan dengan Go dan Gin
// @termsOfService http://swagger.io/terms/

// @contact.name Support Team
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

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
	checkoutHandler := handler.NewCheckoutHandler(checkoutRepo, bookRepo, db)

	authHandler := handler.NewAuthHandler(userRepo)

	config.InitRedis()

	r := gin.Default()

	userRoutes := r.Group("/users")
	{
		// POST /users
		// userRoutes.POST("", userHandler.CreateUser)

		// GET /users
		userRoutes.GET("", userHandler.GetAllUsers)

		// GET /users/:id
		userRoutes.GET("/:id", userHandler.GetUserByID)
	}

	authRoutes := r.Group("/auth")
	{
		// POST /auth/register
		authRoutes.POST("/register", authHandler.Register)

		// POST /auth/login
		authRoutes.POST("/login", authHandler.Login)
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

	checkoutRoutes := r.Group("/checkouts").Use(middleware.AuthMiddleware())
	{
		// POST /checkouts
		checkoutRoutes.POST("", checkoutHandler.CheckoutBook)

		// PUT /checkouts/5/return
		checkoutRoutes.PUT("/:id/return", checkoutHandler.ReturnBook)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Starting server on :8081")
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
