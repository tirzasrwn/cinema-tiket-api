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
	"github.com/stretchr/testify/assert"
)

func TestCreateScreening_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	config.DB = db

	// Mock movie duration query
	movieRows := sqlmock.NewRows([]string{"duration"}).AddRow(120)
	mock.ExpectQuery("SELECT duration FROM movies WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(movieRows)

	// Mock hall capacity query
	hallRows := sqlmock.NewRows([]string{"capacity"}).AddRow(150)
	mock.ExpectQuery("SELECT capacity FROM halls WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(hallRows)

	// Mock insert screening
	mock.ExpectQuery("INSERT INTO screenings").
		WithArgs(1, 1, 1, sqlmock.AnyArg(), sqlmock.AnyArg(), 50000.0, 75000.0, 150, true, true).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	router := setupTestRouter()
	router.POST("/screenings", CreateScreening)

	screeningReq := models.CreateScreeningRequest{
		MovieID:   1,
		TheaterID: 1,
		HallID:    1,
		ShowTime:  time.Now().Add(24 * time.Hour),
		Price:     50000.0,
		Price3D:   75000.0,
		Is3D:      true,
	}

	body, _ := json.Marshal(screeningReq)
	req, _ := http.NewRequest("POST", "/screenings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.True(t, response.Success)
	assert.Equal(t, "Screening created successfully", response.Message)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGetScreening_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	config.DB = db

	now := time.Now()
	showTime := now.Add(24 * time.Hour)

	// Mock screening query
	rows := sqlmock.NewRows([]string{
		"id", "movie_id", "theater_id", "hall_id", "show_time", "end_time",
		"price", "price_3d", "available_seats", "is_3d", "is_available", "created_at", "updated_at",
	}).AddRow(
		1, 1, 1, 1, showTime, showTime.Add(2*time.Hour),
		50000.0, 75000.0, 150, true, true, now, now,
	)

	mock.ExpectQuery("SELECT id, movie_id, theater_id, hall_id, show_time, end_time, price, price_3d, available_seats, is_3d, is_available, created_at, updated_at FROM screenings WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(rows)

	router := setupTestRouter()
	router.GET("/screenings/:id", GetScreening)

	req, _ := http.NewRequest("GET", "/screenings/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.True(t, response.Success)
	assert.Equal(t, "Screening fetched successfully", response.Message)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGetScreening_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	config.DB = db

	// Mock screening not found
	mock.ExpectQuery("SELECT id, movie_id, theater_id, hall_id, show_time, end_time, price, price_3d, available_seats, is_3d, is_available, created_at, updated_at FROM screenings WHERE id = \\$1").
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	router := setupTestRouter()
	router.GET("/screenings/:id", GetScreening)

	req, _ := http.NewRequest("GET", "/screenings/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response models.Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.False(t, response.Success)
	assert.Equal(t, "Screening not found", response.Message)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}
