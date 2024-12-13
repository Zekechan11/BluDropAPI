package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Staff struct {
	ID				int    `json:"id"`
	Staff_name      string `json:"staff_name"`
	Address 		string `json:"address"`
}

func StaffRoutes(r *gin.Engine, db *sqlx.DB) {
	r.GET("/api/get_staff", func(ctx *gin.Context) {
		var staff []Staff
		
		err := db.Select(&staff, "SELECT * FROM staffs")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, staff)
	})

	r.POST("/api/save_staff", func(ctx *gin.Context) {
		var insertStaff struct {
			Staff_name      string `json:"staff_name"`
			Address 		string `json:"address"`
		}

		if err := ctx.ShouldBindJSON(&insertStaff); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		insertQuery := `
    		INSERT INTO staffs (staff_name, address) 
    		VALUES (:staff_name, :address)`

		_, err := db.NamedExec(insertQuery, insertStaff)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save staff: " + err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"Staff": insertStaff})
	})

	r.PUT("/api/update_staff/:id", func(ctx *gin.Context) {
		var updateStaff struct {
			ID				int    `json:"id"`
			Staff_name      string `json:"staff_name"`
			Address 		string `json:"address"`
		}
	
		if err := ctx.ShouldBindJSON(&updateStaff); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id := ctx.Param("id")
    	updateStaff.ID, _ = strconv.Atoi(id)
	
		updateQuery := `
			UPDATE staffs 
			SET staff_name = :staff_name, address = :address
			WHERE ID = :id`
	
		_, err := db.NamedExec(updateQuery, updateStaff)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update staff: " + err.Error()})
			return
		}
	
		ctx.JSON(http.StatusOK, gin.H{"message": "Staff updated successfully", "staff": updateStaff})
	})

	r.DELETE("/api/delete_staff/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
	
		_, err := db.Exec("DELETE FROM Staffs WHERE id = ?", id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete staff: " + err.Error()})
			return
		}
	
		ctx.JSON(http.StatusOK, gin.H{"message": "Staff deleted successfully"})
	})
	
}
