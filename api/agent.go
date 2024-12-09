package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// Updated Agent struct to include area details
type Agent struct {
	ID        int    `json:"id" db:"id"`
	AreaID    int    `json:"area_id" db:"area_id"` // Add db tag
	AgentName string `json:"agent_name" db:"agent_name"` // Add db tag
	AreaName  string `json:"area_name"` // Optional, for joins
}

type InsertAgent struct {
    ID        int    `db:"id"`
    AgentName string `json:"agent_name" db:"agent_name"`  // This must match the column name in your query
    AreaName  string `json:"area_name" db:"area_name"`        // This must match the column name in your query
}


// AgentHandler holds the database connection
type AgentHandler struct {
	DB *sqlx.DB
}

// NewAgentHandler creates a new AgentHandler
func NewAgentHandler(db *sqlx.DB) *AgentHandler {
	return &AgentHandler{DB: db}
}

// CreateAgent handles the creation of a new agent
func (h *AgentHandler) CreateAgent(c *gin.Context) {
	var agent Agent
	if err := c.ShouldBindJSON(&agent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert agent into the database
	query := `INSERT INTO agents (area_id, agent_name) VALUES (:area_id, :agent_name)`
	result, err := h.DB.NamedExec(query, agent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the last inserted ID and set it in the agent struct
	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	agent.ID = int(id)
	c.JSON(http.StatusCreated, agent)
}

func (h *AgentHandler) GetAllAgents(c *gin.Context) {
	query := `
		SELECT a.id, a.agent_name, COALESCE(ar.area, '') AS area_name
FROM agents a 
LEFT JOIN areas ar ON a.area_id = ar.id
	`

	var agents []InsertAgent
	err := h.DB.Select(&agents, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send the list of agents as the response
	c.JSON(http.StatusOK, agents)
}


// UpdateAgent handles the update of an existing agent
func (h *AgentHandler) UpdateAgent(c *gin.Context) {
	id := c.Param("id")
	var agent Agent
	if err := c.ShouldBindJSON(&agent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the agent in the database
	query := `UPDATE agents SET area_id = :area_id, agent_name = :agent_name WHERE id = :id`
	_, err := h.DB.NamedExec(query, map[string]interface{}{
		"id":         id,
		"area_id":    agent.AreaID,
		"agent_name": agent.AgentName,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Agent updated successfully"})
}

// DeleteAgent handles the deletion of an agent
func (h *AgentHandler) DeleteAgent(c *gin.Context) {
	id := c.Param("id")

	// Delete the agent from the database
	query := `DELETE FROM agents WHERE id = :id`
	_, err := h.DB.NamedExec(query, map[string]interface{}{"id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Agent deleted successfully"})
}

// RegisterAgentRoutes registers the Agent routes with the given router
func RegisterAgentRoutes(r *gin.Engine, db *sqlx.DB) {
	agentHandler := NewAgentHandler(db)

	r.POST("/agent", agentHandler.CreateAgent)
	r.GET("/agent", agentHandler.GetAllAgents)
	r.PUT("/agent/:id", agentHandler.UpdateAgent)
	r.DELETE("/agent/:id", agentHandler.DeleteAgent)
}
