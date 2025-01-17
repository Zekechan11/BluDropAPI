package api

import (
	"database/sql"
	"net/http"
	"waterfalls/dto"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// Updated Agent struct to include area details
type Agent struct {
	ID        int    `json:"id" db:"id"`
	AreaID    int    `json:"area_id" db:"area_id"`       // Add db tag
	AgentName string `json:"agent_name" db:"agent_name"` // Add db tag
	AreaName  string `json:"area_name"`                  // Optional, for joins
}

type InsertAgent struct {
	ID        int    `db:"id"`
	AgentName string `json:"agent_name" db:"agent_name"` // This must match the column name in your query
	AreaName  string `json:"area_name" db:"area_name"`   // This must match the column name in your query
}

type DashboardCount struct {
	Payment         float64 `db:"SUM(payment)"`
	NumGallonsOrder int     `db:"SUM(num_gallons_order)"`
	ReturnedGallons int     `db:"SUM(returned_gallons)"`
}

// RegisterAgentRoutes registers the Agent routes with the given router
func AgentRoutes(r *gin.Engine, db *sqlx.DB) {

	r.GET("/v2/api/agent/assigned/:area_id", func(ctx *gin.Context) {
		area_id := ctx.Param("area_id")

		query := `
			SELECT
				CONCAT(firstname, ' ', lastname) AS fullname
			FROM account_staffs
			WHERE area_id = ? AND role = 'Agent'
			LIMIT 1
		`
		var fullname string

		err := db.Get(&fullname, query, area_id)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "No agent found for the specified area"})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch agent: " + err.Error()})
			}
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": fullname})
	})

	r.GET("/v2/api/agent/dashboard/:area_id", func(ctx *gin.Context) {
		area_id := ctx.Param("area_id")

		query := `
			SELECT
				count
			FROM fgs
			WHERE area_id = ?
			LIMIT 1
		`
		var count string
		err := db.Get(&count, query, area_id)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "No agent found for the specified area"})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch agent: " + err.Error()})
			}
			return
		}

		queryDashboardCount := `
			SELECT
				SUM(payment),
				SUM(num_gallons_order),
				SUM(returned_gallons)
			FROM customer_order
			WHERE area_id = ? AND status = 'Completed'
		`

		var dashboardCount DashboardCount
		err = db.Get(&dashboardCount, queryDashboardCount, area_id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"fgs":               count,
			"collected_ammount": dashboardCount.Payment,
			"gallons_delivered": dashboardCount.NumGallonsOrder,
			"gallons_returned":  dashboardCount.ReturnedGallons,
		})
	})

	// DEAD ROUTES ------------------------------------>>

	r.GET("/api/agents", func(ctx *gin.Context) {
		query := `
			SELECT id, firstname, lastname, email, area FROM accounts
			WHERE role = "Staff"
		`
		var agents []dto.AgentsEntity
		err := db.Select(&agents, query)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, agents)
	})
	r.POST("/api/agent", func(ctx *gin.Context) {
		var agent dto.AgentsModel
		if err := ctx.ShouldBindJSON(&agent); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if agent.Role == "" {
			agent.Role = "Staff"
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
			"message": "Agent account created successfully",
			"agentId": lastInsertedID,
		})
	})

	r.POST("/agent", func(ctx *gin.Context) {
		var agent Agent
		if err := ctx.ShouldBindJSON(&agent); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Insert agent into the database
		query := `INSERT INTO agents (area_id, agent_name) VALUES (:area_id, :agent_name)`
		result, err := db.NamedExec(query, agent)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Retrieve the last inserted ID and set it in the agent struct
		id, err := result.LastInsertId()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		agent.ID = int(id)
		ctx.JSON(http.StatusCreated, agent)
	})

	r.GET("/agent", func(ctx *gin.Context) {
		query := `
		SELECT a.id, a.agent_name, COALESCE(ar.area, '') AS area_name
		FROM agents a 
		LEFT JOIN areas ar ON a.area_id = ar.id`

		var agents []InsertAgent
		err := db.Select(&agents, query)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Send the list of agents as the response
		ctx.JSON(http.StatusOK, agents)
	})

	r.PUT("/agent/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var agent Agent
		if err := ctx.ShouldBindJSON(&agent); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update the agent in the database
		query := `UPDATE agents SET area_id = :area_id, agent_name = :agent_name WHERE id = :id`
		_, err := db.NamedExec(query, map[string]interface{}{
			"id":         id,
			"area_id":    agent.AreaID,
			"agent_name": agent.AgentName,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Agent updated successfully"})
	})
	r.DELETE("/agent/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		// Delete the agent from the database
		query := `DELETE FROM agents WHERE id = :id`
		_, err := db.NamedExec(query, map[string]interface{}{"id": id})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Agent deleted successfully"})
	})
}
