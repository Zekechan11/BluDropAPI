package api

import (
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

		query := `INSERT INTO pricing (value, type) VALUES (:value, :type)`
		_, err := db.NamedExec(query, insertPricing)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, insertPricing)
	})

	r.GET("/api/price/:type", func(ctx *gin.Context) {
		utype := ctx.Param("type")
		var price dto.PricingModel

		err := db.Get(&price, "SELECT * FROM pricing WHERE type = ?", utype)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, price)
	})
}
