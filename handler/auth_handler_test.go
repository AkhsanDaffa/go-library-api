package handler

import (
	"bytes"
	"errors"
	"go-library-api/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- 1. MEMBUAT MOCK OBJECT (Pemeran Pengganti) ---
type MockUserRepo struct {
	mock.Mock
}

// Kita pura-pura implementasi fungsi CreateUser
func (m *MockUserRepo) CreateUser(username, email, passwordHash string) (models.User, error) {
	// Rekam argumen yang diterima
	args := m.Called(username, email, passwordHash)

	// Kembalikan data palsu sesuai skenario tes (return index 0 dan 1)
	return args.Get(0).(models.User), args.Error(1)
}

// Kita pura-pura implementasi GetUserByEmail (biar memenuhi interface)
func (m *MockUserRepo) GetUserByEmail(email string) (models.User, error) {
	args := m.Called(email)
	return args.Get(0).(models.User), args.Error(1)
}

// --- 2. MULAI TESTING ---
func TestRegister_Success(t *testing.T) {
	// A. Setup
	gin.SetMode(gin.TestMode) // Mode tes biar log gak berisik
	mockRepo := new(MockUserRepo)
	authHandler := NewAuthHandler(mockRepo) // Inject Repo Palsu ke Handler

	// B. Skenario: "Kalau ada yang minta CreateUser, sukseskan saja!"
	userBaru := models.User{ID: 1, Username: "testuser", Email: "test@example.com"}

	// Ini bagian magic-nya: Kita mendikte Mock apa yang harus dilakukan
	mockRepo.On("CreateUser", "testuser", "test@example.com", mock.Anything).Return(userBaru, nil)

	// C. Bikin Request Pura-pura
	router := gin.Default()
	router.POST("/auth/register", authHandler.Register)

	inputJSON := `{"username":"testuser", "email":"test@example.com", "password":"password123"}`
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBufferString(inputJSON))
	w := httptest.NewRecorder() // Perekam respon (pengganti Postman)

	// D. Eksekusi
	router.ServeHTTP(w, req)

	// E. Validasi (Assert)
	// Pastikan statusnya 201 Created
	assert.Equal(t, http.StatusCreated, w.Code)

	// Pastikan MockRepo benar-benar dipanggil
	mockRepo.AssertExpectations(t)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	// A. Setup
	mockRepo := new(MockUserRepo)
	authHandler := NewAuthHandler(mockRepo)

	// B. Skenario: "Kalau CreateUser dipanggil, ERROR-kan saja!"
	// Kita simulasikan error database
	mockRepo.On("CreateUser", "duplikat", "duplikat@example.com", mock.Anything).
		Return(models.User{}, errors.New("Email already exists"))

	// C. Request
	router := gin.Default()
	router.POST("/auth/register", authHandler.Register)

	inputJSON := `{"username":"duplikat", "email":"duplikat@example.com", "password":"password123"}`
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBufferString(inputJSON))
	w := httptest.NewRecorder()

	// D. Eksekusi
	router.ServeHTTP(w, req)

	// E. Validasi
	// Harusnya error 409 Conflict (sesuai kode handler Anda)
	assert.Equal(t, http.StatusConflict, w.Code)
}
