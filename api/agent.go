package api

import (
	"database/sql"
	"net/http"
	"bludrop-api/dto"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

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
		var count *string
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

		var dashboardCount dto.DashboardCount
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
}
