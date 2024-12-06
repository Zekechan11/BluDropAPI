package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Area struct {
	ID   int    `json:"id"`
	Area string `json:"area"`
}

type AreaHandler struct {
	DB *sqlx.DB
}

func NewAreaHandler(db *sqlx.DB) *AreaHandler {
	return &AreaHandler{DB: db}
}

func (h *AreaHandler) CreateArea(c *gin.Context) {
	var area Area
	if err := c.ShouldBindJSON(&area); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO areas (area) VALUES (:area)`
	result, err := h.DB.NamedExec(query, area)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	area.ID = int(id)
	c.JSON(http.StatusCreated, area)
}

func (h *AreaHandler) GetAllAreas(c *gin.Context) {
	var areas []Area
	err := h.DB.Select(&areas, "SELECT id, area FROM areas")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, areas)
}

func (h *AreaHandler) UpdateArea(c *gin.Context) {
	id := c.Param("id")
	var area Area
	if err := c.ShouldBindJSON(&area); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `UPDATE areas SET area = :area WHERE id = :id`
	_, err := h.DB.NamedExec(query, map[string]interface{}{
		"id":   id,
		"area": area.Area,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Area updated successfully"})
}

func (h *AreaHandler) DeleteArea(c *gin.Context) {
	id := c.Param("id")

	query := `DELETE FROM areas WHERE id = :id`
	_, err := h.DB.NamedExec(query, map[string]interface{}{"id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Area deleted successfully"})
}

func RegisterRoutes(r *gin.Engine, db *sqlx.DB) {
	areaHandler := NewAreaHandler(db)

	r.POST("/area", areaHandler.CreateArea)
	r.GET("/area", areaHandler.GetAllAreas)
	r.PUT("/area/:id", areaHandler.UpdateArea)
	r.DELETE("/area/:id", areaHandler.DeleteArea)
}
