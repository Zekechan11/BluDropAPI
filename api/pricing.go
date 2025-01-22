package api

import (
	"fmt"
	"net/http"
	"waterfalls/dto"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func PricingRoutes(r *gin.Engine, db *sqlx.DB) {
	r.POST("/api/update/price", func(ctx *gin.Context) {
		var insertPricing dto.InsertPricing

		if err := ctx.ShouldBindJSON(&insertPricing); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := `INSERT INTO pricing (dealer, regular) VALUES (:dealer, :regular)`
		_, err := db.NamedExec(query, insertPricing)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, insertPricing)
	})

	r.PUT("/api/price/update", func(ctx *gin.Context) {
		var insertPricing dto.InsertPricing
		
		if err := ctx.ShouldBindJSON(&insertPricing); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := `UPDATE pricing SET dealer = ?, regular = ? WHERE pricing_id = 1`

		_, err := db.Exec(query, insertPricing.Dealer, insertPricing.Regular)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, insertPricing)

	})

	r.GET("/api/price/:type", func(ctx *gin.Context) {
		utype := ctx.Param("type")

		var price float64
		query := fmt.Sprintf("SELECT %s FROM pricing", utype)
			err := db.Get(&price, query)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusOK, price)
	})
}
