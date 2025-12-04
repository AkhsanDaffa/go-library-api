package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go-library-api/config"
	"go-library-api/models"
	"go-library-api/repository"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type BookHandler struct {
	Repo *repository.BookRepository
}

func NewBookHandler(repo *repository.BookRepository) *BookHandler {
	return &BookHandler{Repo: repo}
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	var req struct {
		Title  string `json:"title"`
		Author string `json:"author"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	book, err := h.Repo.CreateBook(req.Title, req.Author, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	c.JSON(http.StatusCreated, book)
}

// func (h *BookHandler) GetAllBooks(c *gin.Context) {
// 	books, err := h.Repo.GetAllBooks()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve books"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, books)
// }

func (h *BookHandler) GetAllBooks(c *gin.Context) {
	ctx := context.Background()

	cacheKey := "books:all"

	val, err := config.RedisClient.Get(ctx, cacheKey).Result()

	if err == nil {
		var books []models.Book
		err = json.Unmarshal([]byte(val), &books)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"source": "redis",
				"data":   books,
			})
			return
		}
	} else if err != redis.Nil {
		fmt.Println("Redis error:", err)
	}

	books, err := h.Repo.GetAllBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	jsonBytes, err := json.Marshal(books)
	if err == nil {
		config.RedisClient.Set(ctx, cacheKey, jsonBytes, 1*time.Minute)
	}

	c.JSON(http.StatusOK, gin.H{
		"source": "database",
		"data":   books,
	})

}

func (h *BookHandler) GetBookByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	book, err := h.Repo.GetBookByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve book"})
		return
	}

	if book.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, book)
}
