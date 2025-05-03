package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func TransactionRoutes(r *gin.Engine, db *sqlx.DB) {
	r.GET("/api/get_transaction", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message":   "Okay",
		})
	})
}