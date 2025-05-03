package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Remittance struct {
	ID              int     `json:"id" db:"id"`
	Date            string  `json:"date" db:"date"`
	AgentID         int     `json:"agent_id" db:"agent_id"`
	AreaID          int     `json:"area_id" db:"area_id"`
	GallonsLoaded   int     `json:"gallons_loaded" db:"gallons_loaded"`
	GallonsSold     int     `json:"gallons_sold" db:"gallons_sold"`
	GallonsCredited int     `json:"gallons_credited" db:"gallons_credited"`
	EmptyReturns    int     `json:"empty_returns" db:"empty_returns"`
	LoanPayments    float64 `json:"loan_payments" db:"loan_payments"`
	NewLoans        float64 `json:"new_loans" db:"new_loans"`
	AmountCollected float64 `json:"amount_collected" db:"amount_collected"`
	ExpectedAmount  float64 `json:"expected_amount" db:"expected_amount"`
	Status          string  `json:"status" db:"status"`
	AgentName       string  `json:"agent_name" db:"-"`
	Area            string  `json:"area" db:"-"`
}

type RemittanceWithDetails struct {
	Remittance
	FirstName string `json:"first_name" db:"firstname"` // From account_staffs
	LastName  string `json:"last_name" db:"lastname"`   // From account_staffs
	AreaName  string `json:"area_name" db:"area"`       // From areas
}

func RegisterRemittanceRoutes(r *gin.Engine, db *sqlx.DB) {
	// Create a new remittance
	r.POST("/v2/api/create_remittance", func(ctx *gin.Context) {
		var remittance Remittance
		if err := ctx.ShouldBindJSON(&remittance); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := `INSERT INTO remittances (
			date, agent_id, area_id, gallons_loaded, gallons_sold, 
			gallons_credited, empty_returns, loan_payments, 
			new_loans, amount_collected, expected_amount, status
		) VALUES (
			:date, :agent_id, :area_id, :gallons_loaded, :gallons_sold, 
			:gallons_credited, :empty_returns, :loan_payments, 
			:new_loans, :amount_collected, :expected_amount, :status
		)`

		result, err := db.NamedExec(query, remittance)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		id, err := result.LastInsertId()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		remittance.ID = int(id)
		ctx.JSON(http.StatusCreated, remittance)
	})

	// Get all remittances with agent and area details
	r.GET("/v2/api/get_remittances", func(ctx *gin.Context) {
		var remittances []RemittanceWithDetails
	
		query := `
			SELECT 
				r.*,
				s.firstname,
				s.lastname,
				a.area
			FROM 
				remittances r
			JOIN 
				account_staffs s ON r.agent_id = s.staff_id
			JOIN 
				areas a ON r.area_id = a.id
			ORDER BY 
				r.date DESC`
	
		err := db.Select(&remittances, query)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	
		// Transform the data to include full agent name
		result := make([]map[string]interface{}, len(remittances))
		for i, rem := range remittances {
			result[i] = map[string]interface{}{
				"id":               rem.ID,
				"date":             rem.Date,
				"agent_id":        rem.AgentID,
				"agent_name":      fmt.Sprintf("%s %s", rem.FirstName, rem.LastName),
				"area_id":         rem.AreaID,
				"area_name":       rem.AreaName,
				"gallons_loaded":   rem.GallonsLoaded,
				"gallons_sold":     rem.GallonsSold,
				"gallons_credited": rem.GallonsCredited,
				"empty_returns":    rem.EmptyReturns,
				"loan_payments":    rem.LoanPayments,
				"new_loans":       rem.NewLoans,
				"amount_collected": rem.AmountCollected,
				"expected_amount": rem.ExpectedAmount,
				"status":          rem.Status,
			}
		}
	
		ctx.JSON(http.StatusOK, result)
	})

	// Filter remittances by date range
	r.GET("/v2/api/get_remittances_by_date", func(ctx *gin.Context) {
		startDate := ctx.Query("start_date")
		endDate := ctx.Query("end_date")

		var remittances []RemittanceWithDetails

		query := `
		SELECT 
			r.*, 
			a.agent_name,
			ar.area
		FROM 
			remittances r
		JOIN 
			agents a ON r.agent_id = a.Id
		JOIN 
			areas ar ON r.area_id = ar.id
		WHERE 
			r.date BETWEEN ? AND ?
		ORDER BY 
			r.date DESC`

		err := db.Select(&remittances, query, startDate, endDate)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, remittances)
	})

	// Get remittances for today
	r.GET("/v2/api/get_todays_remittances", func(ctx *gin.Context) {
		today := time.Now().Format("2006-01-02")
	
		var remittances []RemittanceWithDetails
	
		query := `
			SELECT 
				r.*,
				s.firstname,
				s.lastname,
				a.area
			FROM 
				remittances r
			JOIN 
				account_staffs s ON r.agent_id = s.staff_id
			JOIN 
				areas a ON r.area_id = a.id
			WHERE 
				r.date = ?
			ORDER BY 
				r.id DESC`
	
		err := db.Select(&remittances, query, today)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	
		// Transform the data as above
		result := make([]map[string]interface{}, len(remittances))
		for i, rem := range remittances {
			result[i] = map[string]interface{}{
				"id":               rem.ID,
				"date":             rem.Date,
				"agent_id":        rem.AgentID,
				"agent_name":      fmt.Sprintf("%s %s", rem.FirstName, rem.LastName),
				"area_id":         rem.AreaID,
				"area_name":       rem.AreaName,
				// ... other fields
			}
		}
	
		ctx.JSON(http.StatusOK, result)
	})

	// Get a single remittance
	r.GET("/v2/api/get_remittance/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		var remittance RemittanceWithDetails

		query := `
		SELECT 
			r.*, 
			a.agent_name,
			ar.area
		FROM 
			remittances r
		JOIN 
			agents a ON r.agent_id = a.Id
		JOIN 
			areas ar ON r.area_id = ar.id
		WHERE 
			r.id = ?`

		err := db.Get(&remittance, query, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, remittance)
	})

	// Update a remittance
	r.PUT("/v2/api/update_remittance/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var remittance Remittance

		if err := ctx.ShouldBindJSON(&remittance); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := `
		UPDATE remittances 
		SET 
			date = :date,
			agent_id = :agent_id,
			area_id = :area_id,
			gallons_loaded = :gallons_loaded,
			gallons_sold = :gallons_sold,
			gallons_credited = :gallons_credited,
			empty_returns = :empty_returns,
			loan_payments = :loan_payments,
			new_loans = :new_loans,
			amount_collected = :amount_collected,
			expected_amount = :expected_amount,
			status = :status
		WHERE 
			id = :id`

		_, err := db.NamedExec(query, map[string]interface{}{
			"id":               id,
			"date":             remittance.Date,
			"agent_id":         remittance.AgentID,
			"area_id":          remittance.AreaID,
			"gallons_loaded":   remittance.GallonsLoaded,
			"gallons_sold":     remittance.GallonsSold,
			"gallons_credited": remittance.GallonsCredited,
			"empty_returns":    remittance.EmptyReturns,
			"loan_payments":    remittance.LoanPayments,
			"new_loans":        remittance.NewLoans,
			"amount_collected": remittance.AmountCollected,
			"expected_amount":  remittance.ExpectedAmount,
			"status":           remittance.Status,
		})

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Remittance updated successfully"})
	})

	// Delete a remittance
	r.DELETE("/v2/api/delete_remittance/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		query := `DELETE FROM remittances WHERE id = ?`
		_, err := db.Exec(query, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Remittance deleted successfully"})
	})

	// Filter remittances by status
	r.GET("/v2/api/get_remittances_by_status", func(ctx *gin.Context) {
		status := ctx.Query("status")

		var remittances []RemittanceWithDetails

		query := `
		SELECT 
			r.*, 
			a.agent_name,
			ar.area
		FROM 
			remittances r
		JOIN 
			agents a ON r.agent_id = a.Id
		JOIN 
			areas ar ON r.area_id = ar.id
		WHERE 
			r.status = ?
		ORDER BY 
			r.date DESC`

		err := db.Select(&remittances, query, status)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, remittances)
	})

	// Filter remittances by agent
	r.GET("/v2/api/get_remittances_by_agent", func(ctx *gin.Context) {
		agentId := ctx.Query("agent_id")

		var remittances []RemittanceWithDetails

		query := `
		SELECT 
			r.*, 
			a.agent_name,
			ar.area
		FROM 
			remittances r
		JOIN 
			agents a ON r.agent_id = a.Id
		JOIN 
			areas ar ON r.area_id = ar.id
		WHERE 
			r.agent_id = ?
		ORDER BY 
			r.date DESC`

		err := db.Select(&remittances, query, agentId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, remittances)
	})
}
