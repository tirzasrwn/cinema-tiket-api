package handlers

import (
	"cinema-ticket-api/config"
	"cinema-ticket-api/models"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateScreening godoc
//
//	@Summary		Create a new screening
//	@Description	Create a new movie screening (Admin only)
//	@Tags			screenings
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			screeningRequest	body		models.CreateScreeningRequest			true	"Screening data"
//	@Success		201					{object}	models.Response{data=object{id=int}}	"Screening created successfully"
//	@Failure		400					{object}	models.Response							"Invalid request"
//	@Failure		401					{object}	models.Response							"Unauthorized"
//	@Failure		500					{object}	models.Response							"Internal server error"
//	@Router			/screenings [post]
func CreateScreening(c *gin.Context) {
	var req models.CreateScreeningRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request", err))
		return
	}

	// Calculate end time (show time + movie duration)
	var movieDuration int
	err := config.DB.QueryRow("SELECT duration FROM movies WHERE id = $1", req.MovieID).Scan(&movieDuration)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Movie not found", err))
		return
	}

	endTime := req.ShowTime.Add(time.Duration(movieDuration) * time.Minute)

	// Get hall capacity
	var hallCapacity int
	err = config.DB.QueryRow("SELECT capacity FROM halls WHERE id = $1", req.HallID).Scan(&hallCapacity)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Hall not found", err))
		return
	}

	var screeningID int
	err = config.DB.QueryRow(`
        INSERT INTO screenings 
        (movie_id, theater_id, hall_id, show_time, end_time, price, price_3d, available_seats, is_3d, is_available)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id
    `, req.MovieID, req.TheaterID, req.HallID, req.ShowTime, endTime,
		req.Price, req.Price3D, hallCapacity, req.Is3D, true).Scan(&screeningID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to create screening", err))
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse("Screening created successfully", gin.H{"id": screeningID}))
}

// GetScreenings godoc
//
//	@Summary		Get all screenings
//	@Description	Get list of all available screenings
//	@Tags			screenings
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	models.Response{data=[]models.Screening}	"Screenings fetched successfully"
//	@Failure		401	{object}	models.Response								"Unauthorized"
//	@Failure		500	{object}	models.Response								"Internal server error"
//	@Router			/screenings [get]
func GetScreenings(c *gin.Context) {
	rows, err := config.DB.Query(`
        SELECT id, movie_id, theater_id, hall_id, show_time, end_time, price, 
               price_3d, available_seats, is_3d, is_available, created_at, updated_at
        FROM screenings 
        WHERE show_time > NOW()
        ORDER BY show_time
    `)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to fetch screenings", err))
		return
	}
	defer rows.Close()

	var screenings []models.Screening
	for rows.Next() {
		var s models.Screening
		err := rows.Scan(
			&s.ID, &s.MovieID, &s.TheaterID, &s.HallID, &s.ShowTime, &s.EndTime,
			&s.Price, &s.Price3D, &s.AvailableSeats, &s.Is3D, &s.IsAvailable,
			&s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to scan screening", err))
			return
		}
		screenings = append(screenings, s)
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Screenings fetched successfully", screenings))
}

// GetScreening godoc
//
//	@Summary		Get a specific screening
//	@Description	Get details of a specific screening by ID
//	@Tags			screenings
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int										true	"Screening ID"
//	@Success		200	{object}	models.Response{data=models.Screening}	"Screening fetched successfully"
//	@Failure		400	{object}	models.Response							"Invalid ID"
//	@Failure		401	{object}	models.Response							"Unauthorized"
//	@Failure		404	{object}	models.Response							"Screening not found"
//	@Failure		500	{object}	models.Response							"Internal server error"
//	@Router			/screenings/{id} [get]
func GetScreening(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid screening ID", err))
		return
	}

	var screening models.Screening
	err = config.DB.QueryRow(`
        SELECT id, movie_id, theater_id, hall_id, show_time, end_time, price, 
               price_3d, available_seats, is_3d, is_available, created_at, updated_at
        FROM screenings WHERE id = $1
    `, id).Scan(
		&screening.ID, &screening.MovieID, &screening.TheaterID, &screening.HallID,
		&screening.ShowTime, &screening.EndTime, &screening.Price, &screening.Price3D,
		&screening.AvailableSeats, &screening.Is3D, &screening.IsAvailable,
		&screening.CreatedAt, &screening.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.ErrorResponse("Screening not found", nil))
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse("Database error", err))
		}
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Screening fetched successfully", screening))
}

// UpdateScreening godoc
//
//	@Summary		Update a screening
//	@Description	Update an existing screening (Admin only)
//	@Tags			screenings
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id					path		int								true	"Screening ID"
//	@Param			screeningRequest	body		models.UpdateScreeningRequest	true	"Screening data to update"
//	@Success		200					{object}	models.Response					"Screening updated successfully"
//	@Failure		400					{object}	models.Response					"Invalid request"
//	@Failure		401					{object}	models.Response					"Unauthorized"
//	@Failure		404					{object}	models.Response					"Screening not found"
//	@Failure		500					{object}	models.Response					"Internal server error"
//	@Router			/screenings/{id} [put]
func UpdateScreening(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid screening ID", err))
		return
	}

	var req models.UpdateScreeningRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request", err))
		return
	}

	// Build dynamic update query
	query := "UPDATE screenings SET updated_at = NOW()"
	params := []interface{}{}
	paramCount := 1

	if req.MovieID != 0 {
		query += ", movie_id = $" + strconv.Itoa(paramCount)
		params = append(params, req.MovieID)
		paramCount++
	}

	if req.TheaterID != 0 {
		query += ", theater_id = $" + strconv.Itoa(paramCount)
		params = append(params, req.TheaterID)
		paramCount++
	}

	if req.HallID != 0 {
		query += ", hall_id = $" + strconv.Itoa(paramCount)
		params = append(params, req.HallID)
		paramCount++
	}

	if !req.ShowTime.IsZero() {
		query += ", show_time = $" + strconv.Itoa(paramCount)
		params = append(params, req.ShowTime)
		paramCount++
	}

	if req.Price != 0 {
		query += ", price = $" + strconv.Itoa(paramCount)
		params = append(params, req.Price)
		paramCount++
	}

	if req.Price3D != 0 {
		query += ", price_3d = $" + strconv.Itoa(paramCount)
		params = append(params, req.Price3D)
		paramCount++
	}

	query += ", is_3d = $" + strconv.Itoa(paramCount)
	params = append(params, req.Is3D)
	paramCount++

	query += ", is_available = $" + strconv.Itoa(paramCount)
	params = append(params, req.IsAvailable)
	paramCount++

	query += " WHERE id = $" + strconv.Itoa(paramCount)
	params = append(params, id)

	result, err := config.DB.Exec(query, params...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to update screening", err))
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse("Screening not found", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Screening updated successfully", nil))
}

// DeleteScreening godoc
//
//	@Summary		Delete a screening
//	@Description	Soft delete a screening by setting is_available to false (Admin only)
//	@Tags			screenings
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int				true	"Screening ID"
//	@Success		200	{object}	models.Response	"Screening deleted successfully"
//	@Failure		400	{object}	models.Response	"Invalid ID"
//	@Failure		401	{object}	models.Response	"Unauthorized"
//	@Failure		404	{object}	models.Response	"Screening not found"
//	@Failure		500	{object}	models.Response	"Internal server error"
//	@Router			/screenings/{id} [delete]
func DeleteScreening(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid screening ID", err))
		return
	}

	result, err := config.DB.Exec("UPDATE screenings SET is_available = false WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to delete screening", err))
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse("Screening not found", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Screening deleted successfully", nil))
}
