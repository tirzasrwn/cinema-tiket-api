package models

// LoginRequest represents login credentials
//
//	@Description	User login credentials
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"admin@cinema.com"`
	Password string `json:"password" binding:"required" example:"password"`
}
