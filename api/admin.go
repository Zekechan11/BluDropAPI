package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func AdminRoutes(r *gin.Engine, db *sqlx.DB) {
	r.GET("/api/admin/dashboard", func(ctx *gin.Context) {
		query := `
			SELECT SUM(payment)
			FROM customer_order
		`
		var payment float64
		err := db.Get(&payment, query)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"total_sales": payment,
		})
	})
}
