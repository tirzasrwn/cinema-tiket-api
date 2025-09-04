package handlers

import (
	"cinema-ticket-api/config"
	"cinema-ticket-api/models"
	"cinema-ticket-api/utils"
	"database/sql"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

// Login godoc
//
//	@Summary		User login
//	@Description	Authenticate user and return JWT token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			loginRequest	body		models.LoginRequest							true	"Login credentials"
//	@Success		200				{object}	models.Response{data=models.LoginResponse}	"Login successful"
//	@Failure		400				{object}	models.Response								"Invalid request"
//	@Failure		401				{object}	models.Response								"Invalid credentials"
//	@Failure		500				{object}	models.Response								"Internal server error"
//	@Router			/login [post]
func Login(c *gin.Context) {
	var loginReq models.LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request", err))
		return
	}

	var user models.User
	err := config.DB.QueryRow(`
        SELECT id, email, password_hash, full_name, phone_number, date_of_birth, created_at, updated_at
        FROM users WHERE email = $1 AND email_verified = true
    `, loginReq.Email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName,
		&user.PhoneNumber, &user.DateOfBirth, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// Gunakan error response tanpa passing nil error
			c.JSON(http.StatusUnauthorized, models.Response{
				Success: false,
				Message: "Invalid email or password",
				Error:   "Invalid credentials",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse("Database error", err))
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginReq.Password))
	if err != nil {
		// Gunakan error response tanpa passing nil error
		c.JSON(http.StatusUnauthorized, models.Response{
			Success: false,
			Message: "Invalid email or password",
			Error:   "Invalid credentials",
		})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to generate token", err))
		return
	}

	// Jangan kembalikan password hash di response
	userResponse := models.User{
		ID:          user.ID,
		Email:       user.Email,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		DateOfBirth: user.DateOfBirth,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	response := models.LoginResponse{
		Token: token,
		User:  userResponse,
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Login successful", response))
}
