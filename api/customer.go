package api

import (
	"database/sql"
	"net/http"
	"bludrop-api/dto"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func CustomerRoutes(r *gin.Engine, db *sqlx.DB) {

	r.GET("/v2/api/dashboard/:client_id", func(ctx *gin.Context) {
		clientID := ctx.Param("client_id")

		queryContainers := `
			SELECT total_containers_on_loan 
			FROM containers_on_loan
			WHERE customer_id = ?
		`
		var col int
		err := db.Get(&col, queryContainers, clientID)
		if err != nil {
			if err == sql.ErrNoRows {
				col = 0
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		queryLoan := `
			SELECT SUM(total_price - payment) 
			FROM customer_order
			WHERE customer_id = ?
		`
		var totalLoan float64
		err = db.Get(&totalLoan, queryLoan, clientID)
		if err != nil {
			if err == sql.ErrNoRows {
				totalLoan = 0
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"col":  col,
			"loan": totalLoan,
		})
	})

	r.GET("/v2/api/orders/:client_id", func(ctx *gin.Context) {
		client_id := ctx.Param("client_id")
		query := `
		SELECT
			cu.id,
			cu.num_gallons_order,
			cu.returned_gallons,
			cu.date,
      		cu.total_price,
			cu.payment,
			cu.status,
      		cu.date_created,
      		col.total_containers_on_loan
		FROM customer_order cu
		LEFT JOIN containers_on_loan col ON cu.customer_id = col.customer_id
		WHERE cu.customer_id = ?
		ORDER BY cu.date_created DESC
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
