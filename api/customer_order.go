package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type CustomerOrder struct {
	ID                int    `db:"Id"`
	Num_gallons_order int    `db:"num_gallons_order"`
	Date              string `db:"date"`
	Date_created      string `db:"date_created"`
}

func Customer_OrderRoutes(r *gin.Engine, db *sqlx.DB) {
	r.GET("/api/get_order", func(ctx *gin.Context) {
		var order []CustomerOrder

		err := db.Select(&order, "SELECT * FROM customer_order")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, order)
	})

	r.POST("/api/save_order", func(ctx *gin.Context) {
		var insertCustomerOrder struct {
			Num_gallons_order string `json:"num_gallons_order"`
			Date              string `json:"date"`
		}

		if err := ctx.ShouldBindJSON(&insertCustomerOrder); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		insertQuery := `
    		INSERT INTO customer_order (num_gallons_order, date) 
    		VALUES (:num_gallons_order, :date)`

		// Use NamedExec to execute the query
		_, err := db.NamedExec(insertQuery, insertCustomerOrder)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save order: " + err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"Customer_Order": insertCustomerOrder})
	})
}
