package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// Area struct to represent the area entity
type Area struct {
	ID   int    `json:"id"`
	Area string `json:"area"` // Assuming 'area' is the name of the area
}

// AreaHandler holds the database connection
type AreaHandler struct {
	DB *sql.DB
}

// NewAreaHandler creates a new AreaHandler
func NewAreaHandler(db *sql.DB) *AreaHandler {
	return &AreaHandler{DB: db}
}

// CreateArea handles the creation of a new area
func (h *AreaHandler) CreateArea(c *gin.Context) {
	var area Area
	if err := c.ShouldBindJSON(&area); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.DB.Exec("INSERT INTO areas (area) VALUES (?)", area.Area)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	area.ID = int(id)
	c.JSON(http.StatusCreated, area)
}

// GetAllAreas retrieves all areas from the database
func (h *AreaHandler) GetAllAreas(c *gin.Context) {
	rows, err := h.DB.Query("SELECT id, area FROM areas")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var areas []Area
	for rows.Next() {
		var area Area
		if err := rows.Scan(&area.ID, &area.Area); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		areas = append(areas, area)
	}
	c.JSON(http.StatusOK, areas)
}

// UpdateArea handles the update of an existing area
func (h *AreaHandler) UpdateArea(c *gin.Context) {
	id := c.Param("id")
	var area Area
	if err := c.ShouldBindJSON(&area); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.DB.Exec("UPDATE areas SET area = ? WHERE id = ?", area.Area, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Area updated successfully"})
}

// DeleteArea handles the deletion of an area
func (h *AreaHandler) DeleteArea(c *gin.Context) {
	id := c.Param("id")

	_, err := h.DB.Exec("DELETE FROM areas WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Area deleted successfully"})
}

// RegisterRoutes registers the Area routes with the given router
func RegisterRoutes(r *gin.Engine, db *sql.DB) {
	areaHandler := NewAreaHandler(db)

	r.POST("/area", areaHandler.CreateArea)
	r.GET("/area", areaHandler.GetAllAreas)
	r.PUT("/area/:id", areaHandler.UpdateArea)
	r.DELETE("/area/:id", areaHandler.DeleteArea)
}
