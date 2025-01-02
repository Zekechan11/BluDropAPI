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
	r.GET("/api/schedule", func(ctx *gin.Context) {
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
}
