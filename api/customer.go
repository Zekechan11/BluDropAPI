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




	// OLD API - idk if it sill in use ---------------------------------------->


	r.GET("/api/customers", func (ctx *gin.Context){
		query := `
		SELECT id, firstname, lastname, email, area FROM accounts
		WHERE role = "Customer"
	`
	var agents []dto.CustomerEntity
	err := db.Select(&agents, query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, agents)
	})

	r.POST("/api/customer", func (ctx *gin.Context){
		var agent dto.AgentsModel
		if err := ctx.ShouldBindJSON(&agent); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if agent.Role == "" {
			agent.Role = "Customer"
		}

		query := `
				INSERT INTO accounts (firstname, lastname, email, area, password, role)
				VALUES (:firstname, :lastname, :email, :email, :password, :role)
				`
		result, err := db.NamedExec(query, agent)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		lastInsertedID, err := result.RowsAffected()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"message": "Customer account created successfully",
			"customerId": lastInsertedID,
		})
	})
}