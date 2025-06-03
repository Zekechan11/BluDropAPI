package api

import (
	"net/http"
	"bludrop-api/dto"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func AdminRoutes(r *gin.Engine, db *sqlx.DB) {
	r.GET("/api/admin/dashboard", func(ctx *gin.Context) {
		var price dto.PricingModel
		
		query := `
			SELECT SUM(payment)
			FROM customer_order
		`
		var payment float64
		salesErr := db.Get(&payment, query)
		if salesErr != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": salesErr.Error()})
			return
		}

		priceErr := db.Get(&price, "SELECT * FROM pricing LIMIT 1")
			if priceErr != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": priceErr.Error()})
				return
			}

		ctx.JSON(http.StatusOK, gin.H{
			"total_sales": payment,
			"pricing": gin.H{
				"dealer":  price.Dealer,
				"regular": price.Regular,
			},
		})
	})
}
