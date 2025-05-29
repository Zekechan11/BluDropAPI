package api

import (
	"fmt"
	"net/http"
	"waterfalls/dto"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func ManageProfileRoutes(r *gin.Engine, db *sqlx.DB) {
	r.POST("/api/profile/edit/:type", func(ctx *gin.Context) {
		updateType := ctx.Param("type")

		tableConfig := map[string]struct {
			TableName string
			IDField   string
		}{
			"management": {TableName: "account_staffs", IDField: "staff_id"},
			"customer":   {TableName: "account_clients", IDField: "client_id"},
		}

		config, ok := tableConfig[updateType]
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile type"})
			return
		}

		var profile dto.Profile
		if err := ctx.ShouldBindJSON(&profile); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := fmt.Sprintf(`
			UPDATE %s
			SET
				firstname = :firstname,
				lastname = :lastname,
				email = :email
			WHERE %s = :id`, config.TableName, config.IDField)

		_, err := db.NamedExec(query, profile)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
	})

	r.POST("/api/profile/change-password/:type", func(ctx *gin.Context) {
		updateType := ctx.Param("type")

		tableConfig := map[string]struct {
			TableName string
			IDField   string
		}{
			"management": {TableName: "account_staffs", IDField: "staff_id"},
			"customer":   {TableName: "account_clients", IDField: "client_id"},
		}

		config, ok := tableConfig[updateType]
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile type"})
			return
		}

		var password dto.Password
		if err := ctx.ShouldBindJSON(&password); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		selectQuery := fmt.Sprintf(`SELECT password FROM %s WHERE %s = ?`, config.TableName, config.IDField)
		var oldPassword string
		err := db.Get(&oldPassword, selectQuery, password.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}

		if password.CurrentPassword != oldPassword {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Current password is incorrect"})
			return
		}

		if password.NewPassword == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "New password cannot be empty"})
			return
		}

		updateQuery := fmt.Sprintf(`
			UPDATE %s
			SET password = :password
			WHERE %s = :id`,
		config.TableName, config.IDField)

		updateData := map[string]any{
			"id":       password.ID,
			"password": password.NewPassword,
		}

		_, err = db.NamedExec(updateQuery, updateData)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
	})
}
