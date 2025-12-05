package handler

import (
	"go-library-api/repository"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// type AuthHandler struct {
// 	UserRepo *repository.UserRepository
// }

// func NewAuthHandler(userRepo *repository.UserRepository) *AuthHandler {
// 	return &AuthHandler{UserRepo: userRepo}
// }

type AuthHandler struct {
	UserRepo repository.UserRepositoryInterface
}

func NewAuthHandler(userRepo repository.UserRepositoryInterface) *AuthHandler {
	return &AuthHandler{UserRepo: userRepo}
}

type RegisterInput struct {
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// @Summary      Register User
// @Description  Mendaftarkan pengguna baru dengan username, email, dan password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body handler.RegisterInput true "Payload Body (JSON)"
// @Success      201  {object}  models.User
// @Failure      400  {object}  map[string]any
// @Failure      409  {object}  map[string]any
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input RegisterInput

	if err := c.BindJSON(&input); err != nil {
		errorMessages := formatValidationError(err)

		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessages})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user, err := h.UserRepo.CreateUser(input.Username, input.Email, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email atau username sudah digunakan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

// @Summary      Login User
// @Description  Menukar email & password dengan Token JWT untuk akses API
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body handler.LoginInput true "Payload Body (JSON)"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]any
// @Failure      401  {object}  map[string]any
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input LoginInput

	// 1. Validasi Input
	if err := c.BindJSON(&input); err != nil {
		errorMessages := formatValidationError(err)

		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessages})
		return
	}

	// 2. Cari user berdasarkan email (menggunakan fungsi repo yang baru)
	user, err := h.UserRepo.GetUserByEmail(input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	// 3. Cek jika user tidak ditemukan
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 4. Bandingkan password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password or email"})
		return
	}

	// 5. Generate JWT Token
	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT secret not set"})
		return
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"usr": user.Username,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

func formatValidationError(err error) []string {
	var errors []string

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				errors = append(errors, e.Field()+" is required")
			case "min":
				errors = append(errors, e.Field()+" must be at least "+e.Param()+" characters long")
			case "email":
				errors = append(errors, e.Field()+" must be a valid email address")
			default:
				errors = append(errors, e.Field()+" is not valid")
			}
		}
	}
	if len(errors) == 0 {
		errors = append(errors, err.Error())
	}

	return errors
}
