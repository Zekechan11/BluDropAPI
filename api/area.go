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

func RegisterRoutes(r *gin.Engine, db *sqlx.DB) {
	r.POST("/area", func(ctx *gin.Context) {
		var area Area
		if err := ctx.ShouldBindJSON(&area); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := `INSERT INTO areas (area) VALUES (:area)`
		result, err := db.NamedExec(query, area)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		id, err := result.LastInsertId()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		area.ID = int(id)
		ctx.JSON(http.StatusCreated, area)
	})

	r.GET("/area", func(ctx *gin.Context) {
		var areas []Area
		err := db.Select(&areas, "SELECT id, area FROM areas")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, areas)
	})

	r.PUT("/area/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var area Area
		if err := ctx.ShouldBindJSON(&area); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := `UPDATE areas SET area = :area WHERE id = :id`
		_, err := db.NamedExec(query, map[string]interface{}{
			"id":   id,
			"area": area.Area,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Area updated successfully"})
	})
	
	r.DELETE("/area/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		query := `DELETE FROM areas WHERE id = :id`
		_, err := db.NamedExec(query, map[string]interface{}{"id": id})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Area deleted successfully"})
	})
}
