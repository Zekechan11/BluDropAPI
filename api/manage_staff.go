package api

import (
	"net/http"
	"strconv"
	"waterfalls/dto"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Staff struct {
	ID         int    `json:"id"`
	Staff_name string `json:"staff_name"`
	Address    string `json:"address"`
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
			Staff_name string `json:"staff_name"`
			Address    string `json:"address"`
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
			ID         int    `json:"id"`
			Staff_name string `json:"staff_name"`
			Address    string `json:"address"`
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

	//
	//
	// NEW ROUTES
	//
	//

	r.GET("/v2/api/get_staff", func(ctx *gin.Context) {
		var staff []dto.StaffModel

		err := db.Select(&staff, "SELECT * FROM staff_accounts")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": staff})
	})

	// TODO: fix can't select one
	r.GET("/v2/api/get_staff/:staff_id", func(ctx *gin.Context) {

		staff_id := ctx.Param("staff_id")

		var staff dto.StaffModel

		staff.StaffId,_ = strconv.Atoi(staff_id)

		err := db.Select(&staff, "SELECT staff_id, firstname FROM staff_accounts WHERE staff_id = :staff_id")
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

		insertStaff.Role = role

		insertQuery := `
    		INSERT INTO staff_accounts (firstname, lastname, email, password, role) 
    		VALUES (:firstname, :lastname, :email, :password, :role)`

		_, err := db.NamedExec(insertQuery, insertStaff)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save staff: " + err.Error()})
			return
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
			UPDATE staff_accounts 
			SET firstname = :firstname, lastname = :lastname, email = :email, password = :password
			WHERE staff_id = :staff_id`

		_, err := db.NamedExec(updateQuery, updateStaff)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update staff: " + err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Staff updated successfully", "data": updateStaff})
	})

	r.DELETE("/v2/api/delete_staff/:staff_id", func(ctx *gin.Context) {
		staff_id := ctx.Param("staff_id")

		_, err := db.Exec("DELETE FROM staff_accounts WHERE staff_id = ?", staff_id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete staff: " + err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Staff deleted successfully"})
	})

}
