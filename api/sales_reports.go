package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// SalesReport represents the structure for sales data by area
type SalesReport struct {
	AreaID               int     `db:"area_id" json:"area_id"`
	AreaName             string  `db:"area_name" json:"area_name"`
	TotalOrders          int     `db:"total_orders" json:"total_orders"`
	TotalGallonsSold     int     `db:"total_gallons_sold" json:"total_gallons_sold"`
	TotalGallonsReturned int     `db:"total_gallons_returned" json:"total_gallons_returned"`
	TotalRevenue         float64 `db:"total_revenue" json:"total_revenue"`
	TotalPayments        float64 `db:"total_payments" json:"total_payments"`
	OutstandingBalance   float64 `db:"outstanding_balance" json:"outstanding_balance"`
	StartDate            string  `json:"start_date"`
	EndDate              string  `json:"end_date"`
}

func SalesReportRoutes(r *gin.Engine, db *sqlx.DB) {
	r.GET("/api/get_sales_by_area", func(ctx *gin.Context) {
		// Get date range parameters
		startDate := ctx.DefaultQuery("start_date", "")
		endDate := ctx.DefaultQuery("end_date", "")

		// Validate date parameters
		if startDate == "" || endDate == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Both start_date and end_date are required (format: YYYY-MM-DD)",
			})
			return
		}

		// Parse dates
		_, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid start_date format. Use YYYY-MM-DD",
			})
			return
		}

		_, err = time.Parse("2006-01-02", endDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid end_date format. Use YYYY-MM-DD",
			})
			return
		}

		// SQL query to get sales data grouped by area
		query := `
			SELECT 
				a.id as area_id,
				a.area as area_name,
				COUNT(co.id) as total_orders,
				SUM(co.num_gallons_order) as total_gallons_sold,
				SUM(co.returned_gallons) as total_gallons_returned,
				SUM(co.total_price) as total_revenue,
				SUM(co.payment) as total_payments,
				SUM(co.total_price - co.payment) as outstanding_balance
			FROM areas a
			LEFT JOIN customer_order co ON a.id = co.area_id
			WHERE co.date >= ? 
			AND co.date <= ?
			AND co.status = 'Completed'
			GROUP BY a.id, a.area
			ORDER BY total_revenue DESC
		`

		var salesReports []SalesReport
		err = db.Select(&salesReports, query, startDate, endDate)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch sales data: " + err.Error(),
			})
			return
		}

		// Add date range to each report
		for i := range salesReports {
			salesReports[i].StartDate = startDate
			salesReports[i].EndDate = endDate
		}

		// Calculate totals across all areas
		var totals = SalesReport{
			AreaID:               0,
			AreaName:             "All Areas",
			TotalOrders:          0,
			TotalGallonsSold:     0,
			TotalGallonsReturned: 0,
			TotalRevenue:         0,
			TotalPayments:        0,
			OutstandingBalance:   0,
			StartDate:            startDate,
			EndDate:              endDate,
		}

		for _, report := range salesReports {
			totals.TotalOrders += report.TotalOrders
			totals.TotalGallonsSold += report.TotalGallonsSold
			totals.TotalGallonsReturned += report.TotalGallonsReturned
			totals.TotalRevenue += report.TotalRevenue
			totals.TotalPayments += report.TotalPayments
			totals.OutstandingBalance += report.OutstandingBalance
		}

		// Return both area-wise breakdown and totals
		ctx.JSON(http.StatusOK, gin.H{
			"area_reports": salesReports,
			"totals":       totals,
		})
	})
}
