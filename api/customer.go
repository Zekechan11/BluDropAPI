package api

import (
	"net/http"
	"waterfalls/dto"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func CustomerRoutes(r *gin.Engine, db *sqlx.DB) {
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