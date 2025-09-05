package handlers

import (
	"bytes"
	"cinema-ticket-api/config"
	"cinema-ticket-api/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestLogin_Success(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	config.DB = db

	// Test data
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	testUser := models.User{
		ID:           1,
		Email:        "admin@cinema.com",
		PasswordHash: string(hashedPassword),
		FullName:     "Admin User",
		PhoneNumber:  "08123456789",
		DateOfBirth:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Mock expectations
	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "full_name", "phone_number", "date_of_birth", "created_at", "updated_at"}).
		AddRow(testUser.ID, testUser.Email, testUser.PasswordHash, testUser.FullName,
			testUser.PhoneNumber, testUser.DateOfBirth, testUser.CreatedAt, testUser.UpdatedAt)

	mock.ExpectQuery("SELECT id, email, password_hash, full_name, phone_number, date_of_birth, created_at, updated_at FROM users WHERE email = \\$1 AND email_verified = true").
		WithArgs("admin@cinema.com").
		WillReturnRows(rows)

	// Setup router and request
	router := setupTestRouter()
	router.POST("/login", Login)

	loginReq := models.LoginRequest{
		Email:    "admin@cinema.com",
		Password: "password",
	}

	body, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.True(t, response.Success)
	assert.Equal(t, "Login successful", response.Message)

	// Verify mock expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	config.DB = db

	// Mock - user not found
	mock.ExpectQuery("SELECT id, email, password_hash, full_name, phone_number, date_of_birth, created_at, updated_at FROM users WHERE email = \\$1 AND email_verified = true").
		WithArgs("nonexistent@cinema.com").
		WillReturnError(sql.ErrNoRows)

	router := setupTestRouter()
	router.POST("/login", Login)

	loginReq := models.LoginRequest{
		Email:    "nonexistent@cinema.com",
		Password: "wrongpassword",
	}

	body, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response models.Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.False(t, response.Success)
	assert.Equal(t, "Invalid email or password", response.Message)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	config.DB = db

	// Test data with correct password hash
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	testUser := models.User{
		ID:           1,
		Email:        "admin@cinema.com",
		PasswordHash: string(hashedPassword),
		FullName:     "Admin User",
		PhoneNumber:  "08123456789",
		DateOfBirth:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "full_name", "phone_number", "date_of_birth", "created_at", "updated_at"}).
		AddRow(testUser.ID, testUser.Email, testUser.PasswordHash, testUser.FullName,
			testUser.PhoneNumber, testUser.DateOfBirth, testUser.CreatedAt, testUser.UpdatedAt)

	mock.ExpectQuery("SELECT id, email, password_hash, full_name, phone_number, date_of_birth, created_at, updated_at FROM users WHERE email = \\$1 AND email_verified = true").
		WithArgs("admin@cinema.com").
		WillReturnRows(rows)

	router := setupTestRouter()
	router.POST("/login", Login)

	loginReq := models.LoginRequest{
		Email:    "admin@cinema.com",
		Password: "wrongpassword", // Wrong password
	}

	body, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response models.Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.False(t, response.Success)
	assert.Equal(t, "Invalid email or password", response.Message)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestLogin_InvalidRequest(t *testing.T) {
	router := setupTestRouter()
	router.POST("/login", Login)

	// Invalid JSON
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.False(t, response.Success)
	assert.Contains(t, response.Message, "Invalid request")
}
