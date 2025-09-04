package models

import "time"

// Response represents a standard API response
//
//	@Description	Standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// LoginResponse represents login response data
//
//	@Description	Login response data
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  User   `json:"user"`
}

// User represents user data
//
//	@Description	User information
type User struct {
	ID           int       `json:"id" example:"1"`
	Email        string    `json:"email" example:"admin@cinema.com"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"full_name" example:"Admin User"`
	PhoneNumber  string    `json:"phone_number" example:"08123456789"`
	DateOfBirth  time.Time `json:"date_of_birth" example:"1990-01-01T00:00:00Z"`
	CreatedAt    time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt    time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

func SuccessResponse(message string, data interface{}) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(message string, err error) Response {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}

	return Response{
		Success: false,
		Message: message,
		Error:   errorMsg,
	}
}
