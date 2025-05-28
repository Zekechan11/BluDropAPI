package api

import (
	"net/http"
	"strconv"
	"strings"
	"waterfalls/dto"
	"waterfalls/util"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func StaffRoutes(r *gin.Engine, db *sqlx.DB) {

	r.GET("/v2/api/get_staff", func(ctx *gin.Context) {
		var staff []dto.StaffModel

		err := db.Select(&staff, "SELECT * FROM account_staffs")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": staff})
	})

	r.GET("/v2/api/get_staff/all/:role", func(ctx *gin.Context) {
		role := ctx.Param("role")

		var staff []dto.StaffModel

		err := db.Select(&staff, "SELECT s.*, a.area FROM account_staffs s LEFT JOIN areas a ON id = area_id WHERE role = ?", role)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, staff)
	})

	// TODO: fix can't select one
	r.GET("/v2/api/get_staff/:staff_id", func(ctx *gin.Context) {

		staff_id := ctx.Param("staff_id")

		id, _ := strconv.Atoi(staff_id)

		var staff dto.StaffModel

		err := db.Get(&staff, "SELECT staff_id, firstname FROM account_staffs WHERE staff_id = ?", id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": staff})
	})

	// role: Admin', 'Staff', 'Agent'
	r.POST("/v2/api/create_staff/:role", func(ctx *gin.Context) {

		role := ctx.Param("role")

		var insertStaff dto.InsertStaff

		if err := ctx.ShouldBindJSON(&insertStaff); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if role == "Agent" {
			insertStaff.Email = strings.ToLower(insertStaff.Email)

			exists, err := util.SatffEmailCheck(db, insertStaff.Email)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check email: " + err.Error()})
				return
			}

			if exists {
				ctx.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
				return
			}
		}

		insertStaff.Role = role

		insertQuery := `
    		INSERT INTO account_staffs (firstname, lastname, email, password, role, area_id) 
    		VALUES (:firstname, :lastname, :email, :password, :role, :area_id)`

		result, err := db.NamedExec(insertQuery, insertStaff)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save staff: " + err.Error()})
			return
		}

		insertStaff.StaffId, err = result.LastInsertId()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve inserted ID"})
			return
		}

		if insertStaff.Role == "Agent" {
			staffOnly := true
			_, err := util.CreateChatID(db, insertStaff.StaffId, insertStaff.AreaId, &staffOnly)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to create or retrieve chat ID",
					"details": err.Error(),
				})
				return
			}
		}

		ctx.JSON(http.StatusOK, gin.H{"data": insertStaff})
	})

	r.PUT("/v2/api/update_staff/:staff_id", func(ctx *gin.Context) {
		
		var updateStaff dto.StaffModel

		if err := ctx.ShouldBindJSON(&updateStaff); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		staff_id := ctx.Param("staff_id")
		updateStaff.StaffId, _ = strconv.Atoi(staff_id)

		updateQuery := `
			UPDATE account_staffs 
			SET
				firstname = :firstname,
				lastname = :lastname,
				email = :email,
				password = :password,
				area_id = :area_id
			WHERE staff_id = :staff_id`

		_, err := db.NamedExec(updateQuery, updateStaff)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update staff: " + err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Staff updated successfully", "data": updateStaff})
	})

	r.PUT("/v2/api/update_staff/area/:staff_id", func(ctx *gin.Context) {
		
		var updateStaff dto.StaffModel

		if err := ctx.ShouldBindJSON(&updateStaff); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		staff_id := ctx.Param("staff_id")
		updateStaff.StaffId, _ = strconv.Atoi(staff_id)

		updateQuery := `
			UPDATE account_staffs 
			SET
				area_id = :area_id
			WHERE staff_id = :staff_id`

		_, err := db.NamedExec(updateQuery, updateStaff)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update staff area: " + err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Staff area updated successfully", "data": updateStaff})
	})

	r.DELETE("/v2/api/delete_staff/:staff_id", func(ctx *gin.Context) {
		staff_id := ctx.Param("staff_id")

		_, err := db.Exec("DELETE FROM account_staffs WHERE staff_id = ?", staff_id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete staff: " + err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Staff deleted successfully"})
	})
}
