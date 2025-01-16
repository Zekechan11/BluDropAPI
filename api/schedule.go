package api

import (
	"net/http"
	"time"

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

func getNextWeekday(current time.Time, weekday time.Weekday) time.Time {
	// Calculate the date of the next occurrence of the specified weekday
	offset := (int(weekday) - int(current.Weekday()) + 7) % 7
	if offset == 0 {
		offset = 7 // ensure it's next week, not today
	}
	return current.AddDate(0, 0, offset)
}

func ScheduleRoutes(r *gin.Engine, db *sqlx.DB) {
	r.GET("/api/get_schedule", func(ctx *gin.Context) {
		var schedule Schedules

		// Replace with actual database query
		err := db.Get(&schedule, "SELECT * FROM schedules")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var availableDays []gin.H
		currentDate := time.Now()

		if schedule.Monday {
			nextMonday := getNextWeekday(currentDate, time.Monday)
			availableDays = append(availableDays, gin.H{"day": "Monday", "date": nextMonday.Format("2006-01-02")})
		}
		if schedule.Tuesday {
			nextTuesday := getNextWeekday(currentDate, time.Tuesday)
			availableDays = append(availableDays, gin.H{"day": "Tuesday", "date": nextTuesday.Format("2006-01-02")})
		}
		if schedule.Wednesday {
			nextWednesday := getNextWeekday(currentDate, time.Wednesday)
			availableDays = append(availableDays, gin.H{"day": "Wednesday", "date": nextWednesday.Format("2006-01-02")})
		}
		if schedule.Thursday {
			nextThursday := getNextWeekday(currentDate, time.Thursday)
			availableDays = append(availableDays, gin.H{"day": "Thursday", "date": nextThursday.Format("2006-01-02")})
		}
		if schedule.Friday {
			nextFriday := getNextWeekday(currentDate, time.Friday)
			availableDays = append(availableDays, gin.H{"day": "Friday", "date": nextFriday.Format("2006-01-02")})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"days":        availableDays,
		})
	})

	r.GET("/api/admin/get_schedule", func(ctx *gin.Context) {
		var schedule Schedules
	
		err := db.Get(&schedule, "SELECT * FROM schedules")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	
		var availableDays []gin.H
		currentDate := time.Now()
	
		{
			nextMonday := getNextWeekday(currentDate, time.Monday)
			availableDays = append(availableDays, gin.H{
				"day":   "Monday",
				"date":  nextMonday.Format("01-02-06"),
				"type":  schedule.Monday,
			})
		}
	
		{
			nextTuesday := getNextWeekday(currentDate, time.Tuesday)
			availableDays = append(availableDays, gin.H{
				"day":   "Tuesday",
				"date":  nextTuesday.Format("01-02-06"),
				"type":  schedule.Tuesday,
			})
		}
	
		{
			nextWednesday := getNextWeekday(currentDate, time.Wednesday)
			availableDays = append(availableDays, gin.H{
				"day":   "Wednesday",
				"date":  nextWednesday.Format("01-02-06"),
				"type":  schedule.Wednesday,
			})
		}
	
		{
			nextThursday := getNextWeekday(currentDate, time.Thursday)
			availableDays = append(availableDays, gin.H{
				"day":   "Thursday",
				"date":  nextThursday.Format("01-02-06"),
				"type":  schedule.Thursday,
			})
		}
	
		{
			nextFriday := getNextWeekday(currentDate, time.Friday)
			availableDays = append(availableDays, gin.H{
				"day":   "Friday",
				"date":  nextFriday.Format("01-02-06"),
				"type":  schedule.Friday,
			})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"days":        availableDays,
		})
	})	

	r.PUT("/api/admin/update_schedule", func(ctx *gin.Context) {
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
			WHERE schedule_id = 1
		`

		params := map[string]interface{}{
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

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Schedule updated successfully",
		})
	})
}
