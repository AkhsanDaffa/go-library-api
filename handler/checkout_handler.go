package handler

import (
	"go-library-api/repository"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CheckoutHandler struct {
	CheckoutRepo *repository.CheckoutRepository
	BookRepo     *repository.BookRepository
}

func NewCheckoutHandler(checkoutRepo *repository.CheckoutRepository, bookRepo *repository.BookRepository) *CheckoutHandler {
	return &CheckoutHandler{
		CheckoutRepo: checkoutRepo,
		BookRepo:     bookRepo,
	}
}

func (h *CheckoutHandler) CheckoutBook(c *gin.Context) {
	var req struct {
		UserID int `json:"user_id"`
		BookID int `json:"book_id"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Check Book Availability

	book, err := h.BookRepo.GetBookByID(req.BookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve book"})
		return
	}

	if book.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	if !book.IsAvailable {
		c.JSON(http.StatusConflict, gin.H{"error": "Book is not available for checkout"})
		return
	}

	// Create Checkout Record
	checkout, err := h.CheckoutRepo.CreateCheckout(req.UserID, req.BookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to checkout book"})
		return
	}

	// Update: Set Satus buku to not available (false)
	if err := h.BookRepo.UpdateBookAvailability(req.BookID, false); err != nil {
		log.Printf("Failed to update status book %d: %v", req.BookID, err)
	}
	c.JSON(http.StatusCreated, checkout)
}

func (h *CheckoutHandler) ReturnBook(c *gin.Context) {
	idStr := c.Param("id")
	checkoutID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID checkout tidak valid"})
		return
	}

	checkout, err := h.CheckoutRepo.GetCheckoutByID(checkoutID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Data checkout tidak ditemukan"})
		return
	}

	if checkout.ReturnDate.Valid {
		c.JSON(http.StatusConflict, gin.H{"error": "Buku sudah dikembalikan"})
		return
	}

	if err := h.CheckoutRepo.UpdateReturnDate(checkoutID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate tanggal pengembalian"})
		return
	}

	if err := h.BookRepo.UpdateBookAvailability(checkout.BookID, true); err != nil {
		log.Printf("Gagal mengupdate status buku %d: %v", checkout.BookID, err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Buku berhasil dikembalikan"})
}
