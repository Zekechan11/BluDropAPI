package api

import (
	"net/http"
	"waterfalls/dto"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func FGSRoutes(r *gin.Engine, db *sqlx.DB) {
	r.GET("/api/fgs/agent", func(ctx *gin.Context) {
		var agent []dto.StaffFGS
		query := `
			SELECT
				s.staff_id,
				CONCAT(s.firstname, ' ', s.lastname) AS fullname,
				a.id,
				a.area,
				f.fgs_id,
				f.count
			FROM account_staffs s
			LEFT JOIN areas a ON a.id = s.area_id
			LEFT JOIN fgs f ON f.area_id = s.area_id
			WHERE s.role = 'Agent'
		`
		err := db.Select(&agent, query)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, agent)
	})

	r.POST("/api/fgs/add", func(ctx *gin.Context) {
		var insertFGS dto.InsertFGS
	
		if err := ctx.ShouldBindJSON(&insertFGS); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	
		if insertFGS.FGGId != 0 {
			var currentCount int
			err := db.Get(&currentCount, "SELECT count FROM fgs WHERE fgs_id = ?", insertFGS.FGGId)
			if err != nil {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "FGS ID not found"})
				return
			}
	
			newCount := currentCount + insertFGS.Count
			updateQuery := `
				UPDATE fgs 
				SET count = :count 
				WHERE fgs_id = :fgs_id
			`
			_, err = db.NamedExec(updateQuery, map[string]interface{}{
				"count":  newCount,
				"fgs_id": insertFGS.FGGId,
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
	
			ctx.JSON(http.StatusOK, gin.H{
				"message": "UPDATE OK",
				"fgs_id": insertFGS.FGGId,
				"new_count": newCount,
			})
		} else {
			query := `
				INSERT INTO fgs (area_id, count)
				VALUES (:area_id, :count)
			`
			result, err := db.NamedExec(query, insertFGS)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
	
			lastInsertedID, err := result.RowsAffected()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
				return
			}
	
			ctx.JSON(http.StatusCreated, gin.H{
				"message": "OK",
				"fgs_id": lastInsertedID,
			})
		}
	})	
}