package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// Updated Agent struct to include area details
type Agent struct {
	ID        int    `json:"id"`
	AreaID    int    `json:"area_id" binding:"required"` // area_id is required
	AgentName string `json:"agent_name" binding:"required"` // agent_name is required
	AreaName  string `json:"area_name"` // New field to hold the area name
}

// AgentHandler holds the database connection
type AgentHandler struct {
	DB *sql.DB
}

// NewAgentHandler creates a new AgentHandler
func NewAgentHandler(db *sql.DB) *AgentHandler {
	return &AgentHandler{DB: db}
}

// CreateAgent handles the creation of a new agent
func (h *AgentHandler) CreateAgent(c *gin.Context) {
	var agent Agent
	if err := c.ShouldBindJSON(&agent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.DB.Exec("INSERT INTO agents (area_id, agent_name) VALUES (?, ?)", agent.AreaID, agent.AgentName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	agent.ID = int(id)
	c.JSON(http.StatusCreated, agent)
}

// GetAllAgents retrieves all agents along with their area details from the database
func (h *AgentHandler) GetAllAgents(c *gin.Context) {
	query := `
		SELECT a.id, a.area_id, a.agent_name, ar.area 
		FROM agents a 
		LEFT JOIN areas ar ON a.area_id = ar.id
	`
	rows, err := h.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var agents []Agent
	for rows.Next() {
		var agent Agent
		if err := rows.Scan(&agent.ID, &agent.AreaID, &agent.AgentName, &agent.AreaName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		agents = append(agents, agent)
	}
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

	_, err := h.DB.Exec("UPDATE agents SET area_id = ?, agent_name = ? WHERE id = ?", agent.AreaID, agent.AgentName, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Agent updated successfully"})
}

// DeleteAgent handles the deletion of an agent
func (h *AgentHandler) DeleteAgent(c *gin.Context) {
	id := c.Param("id")

	_, err := h.DB.Exec("DELETE FROM agents WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Agent deleted successfully"})
}

// RegisterAgentRoutes registers the Agent routes with the given router
func RegisterAgentRoutes(r *gin.Engine, db *sql.DB) {
	agentHandler := NewAgentHandler(db)

	r.POST("/agent", agentHandler.CreateAgent)
	r.GET("/agent", agentHandler.GetAllAgents)
	r.PUT("/agent/:id", agentHandler.UpdateAgent)
	r.DELETE("/agent/:id", agentHandler.DeleteAgent)
}
