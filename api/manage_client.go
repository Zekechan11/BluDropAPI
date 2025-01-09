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

func ClientRoutes(r *gin.Engine, db *sqlx.DB) {
	r.GET("/v2/api/get_client/all/:status", func(ctx *gin.Context) {
		status := ctx.Param("status")

		var client []dto.ClientModel

		query := `
			SELECT c.*, area FROM account_clients c
			LEFT JOIN areas a ON a.id = c.area_id
			WHERE status = ?
		`

		err := db.Select(&client, query, status)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": client})
	})

	r.POST("/v2/api/create_client", func(ctx *gin.Context) {

		var insertClient dto.InsertClient

		if err := ctx.ShouldBindJSON(&insertClient); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		insertClient.Username = strings.ToLower(insertClient.Username)
		insertClient.Email = strings.ToLower(insertClient.Email)

		exists, err := util.ClientUsernameOrEmailCheck(db, insertClient.Username, insertClient.Email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check username/email: " + err.Error()})
			return
		}

		if exists {
			ctx.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
			return
		}

		insertClient.Role = "Customer"

		insertQuery := `
    		INSERT INTO account_clients (firstname, lastname, email, username, password, area_id, role) 
    		VALUES (:firstname, :lastname, :email, :username, :password, :area_id, :role)`

		_, err = db.NamedExec(insertQuery, insertClient)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save client: " + err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": insertClient})
	})

	r.PUT("/v2/api/update_client/:client_id", func(ctx *gin.Context) {
		
		var updateClient dto.ClientModel

		if err := ctx.ShouldBindJSON(&updateClient); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		client_id := ctx.Param("client_id")
		updateClient.ClientId, _ = strconv.Atoi(client_id)

		updateQuery := `
			UPDATE account_clients 
			SET
				firstname = :firstname,
				lastname = :lastname,
				email = :email,
				username = :username,
				password = :password,
				area_id = :area_id
			WHERE client_id = :client_id`

		_, err := db.NamedExec(updateQuery, updateClient)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update client: " + err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Client updated successfully", "data": updateClient})
	})

	r.DELETE("/v2/api/delete_client/:client_id", func(ctx *gin.Context) {
		client_id := ctx.Param("client_id")

		_, err := db.Exec("DELETE FROM account_clients WHERE client_id = ?", client_id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete client: " + err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Client deleted successfully"})
	})
}
