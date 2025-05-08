package api

import (
	"net/http"
	"waterfalls/dto"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func CustomerRoutes(r *gin.Engine, db *sqlx.DB) {

	r.GET("/v2/api/dashboard/:client_id", func(ctx *gin.Context) {
		client_id := ctx.Param("client_id")
	
		queryContainers := `
			SELECT total_containers_on_loan 
			FROM containers_on_loan
			WHERE customer_id = ?
		`
		var col int
		err := db.Get(&col, queryContainers, client_id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	
		queryLoan := `
			SELECT SUM(total_price - payment) 
			FROM customer_order
			WHERE customer_id = ?
		`
	
		var totalLoan float64
		err = db.Get(&totalLoan, queryLoan, client_id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	
		ctx.JSON(http.StatusOK, gin.H{
			"col": col,
			"loan": totalLoan,
		})
	})

	r.GET("/v2/api/orders/:client_id", func (ctx *gin.Context){
	client_id := ctx.Param("client_id")
		query := `
		SELECT
			id,
			num_gallons_order,
			returned_gallons,
			date,
			payment
		FROM customer_order
		WHERE customer_id = ?
	`
	var agents []dto.CustomerTransaction
	err := db.Select(&agents, query, client_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, agents)
	})
}