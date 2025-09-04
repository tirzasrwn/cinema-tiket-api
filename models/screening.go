package models

import (
	"time"
)

// Screening represents a movie screening
//
//	@Description	Movie screening information
type Screening struct {
	ID             int       `json:"id" example:"1"`
	MovieID        int       `json:"movie_id" example:"1"`
	TheaterID      int       `json:"theater_id" example:"1"`
	HallID         int       `json:"hall_id" example:"1"`
	ShowTime       time.Time `json:"show_time" example:"2025-12-25T18:00:00Z"`
	EndTime        time.Time `json:"end_time" example:"2025-12-25T21:00:00Z"`
	Price          float64   `json:"price" example:"50000.00"`
	Price3D        float64   `json:"price_3d" example:"75000.00"`
	AvailableSeats int       `json:"available_seats" example:"150"`
	Is3D           bool      `json:"is_3d" example:"true"`
	IsAvailable    bool      `json:"is_available" example:"true"`
	CreatedAt      time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt      time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// CreateScreeningRequest represents data needed to create a screening
//
//	@Description	Data required to create a new screening
type CreateScreeningRequest struct {
	MovieID   int       `json:"movie_id" binding:"required" example:"1"`
	TheaterID int       `json:"theater_id" binding:"required" example:"1"`
	HallID    int       `json:"hall_id" binding:"required" example:"1"`
	ShowTime  time.Time `json:"show_time" binding:"required" example:"2025-12-25T18:00:00Z"`
	Price     float64   `json:"price" binding:"required" example:"50000.00"`
	Price3D   float64   `json:"price_3d" example:"75000.00"`
	Is3D      bool      `json:"is_3d" example:"true"`
}

// UpdateScreeningRequest represents data needed to update a screening
//
//	@Description	Data required to update an existing screening
type UpdateScreeningRequest struct {
	MovieID     int       `json:"movie_id" example:"1"`
	TheaterID   int       `json:"theater_id" example:"1"`
	HallID      int       `json:"hall_id" example:"1"`
	ShowTime    time.Time `json:"show_time" example:"2025-12-25T18:00:00Z"`
	Price       float64   `json:"price" example:"50000.00"`
	Price3D     float64   `json:"price_3d" example:"75000.00"`
	Is3D        bool      `json:"is_3d" example:"true"`
	IsAvailable bool      `json:"is_available" example:"true"`
}
