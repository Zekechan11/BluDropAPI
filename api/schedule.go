package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Schedules struct {
	ScheduleId int  `json:"schedule_id" db:"schedule_id"`
	Monday     bool `json:"monday" db:"monday"`
	Tuesday    bool `json:"tuesday" db:"tuesday"`
	Wednesday  bool `json:"wednesday" db:"wednesday"`
	Thursday   bool `json:"thursday" db:"thursday"`
	Friday     bool `json:"friday" db:"friday"`
}

func ScheduleRoutes(r *gin.Engine, db *sqlx.DB) {
	r.GET("/api/get_schedule", func(ctx *gin.Context) {
		var schedule Schedules

		err := db.Get(&schedule, "SELECT * FROM schedules")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		availableDays := make(map[string]bool)

		if schedule.Monday {
			availableDays["monday"] = true
		}
		if schedule.Tuesday {
			availableDays["tuesday"] = true
		}
		if schedule.Wednesday {
			availableDays["wednesday"] = true
		}
		if schedule.Thursday {
			availableDays["thursday"] = true
		}
		if schedule.Friday {
			availableDays["friday"] = true
		}

		ctx.JSON(http.StatusOK, gin.H{
			"schedule_id": schedule.ScheduleId,
			"days":        availableDays,
		})
	})

	r.PUT("/api/update_schedule", func(ctx *gin.Context) {
		var req Schedules

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		query := `
			UPDATE schedules
			SET monday = :monday,
			    tuesday = :tuesday,
			    wednesday = :wednesday,
			    thursday = :thursday,
			    friday = :friday
			WHERE schedule_id = :schedule_id
		`

		params := map[string]interface{}{
			"schedule_id": req.ScheduleId,
			"monday":      req.Monday,
			"tuesday":     req.Tuesday,
			"wednesday":   req.Wednesday,
			"thursday":    req.Thursday,
			"friday":      req.Friday,
		}

		_, err := db.NamedExec(query, params)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update schedule"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Schedule updated successfully"})
	})
}
