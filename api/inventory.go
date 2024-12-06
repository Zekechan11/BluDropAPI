package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func Inventory(r *gin.Engine, db *sqlx.DB) {
	r.GET("/api/get_inventory", func(c *gin.Context) {
		var inventory []map[string]interface{}

		err := db.Select(&inventory, "SELECT * FROM Inventory")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, inventory)
	})
}
