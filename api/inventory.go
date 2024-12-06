package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type InventoryItem struct {
	ID     int     `db:"inventory_id"`
	Name   string  `db:"item"`
	NoOfItems  int `db:"no_of_items"`
}

func InventoryRoutes(r *gin.Engine, db *sqlx.DB) {
	r.GET("/api/get_inventory", func(ctx *gin.Context) {
		var inventory []InventoryItem

		err := db.Select(&inventory, "SELECT * FROM Inventory")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, inventory)
	})

	r.POST("/api/save_inventory", func(ctx *gin.Context) {
		var insertInventory struct {
			Name      string `json:"name"`
			NoOfItems int    `json:"no_of_items"`
		}
	
		if err := ctx.ShouldBindJSON(&insertInventory); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	
		insertQuery := `
			INSERT INTO Inventory (Item, noOfItems) 
			VALUES (:Item, :NoOfItems)
			RETURNING inventoryId, Item, noOfItems
		`
	
		var insertedInventory struct {
			InventoryID int    `db:"inventoryId"`
			Item        string `db:"Item"`
			NoOfItems   int    `db:"noOfItems"`
		}
	
		err := db.QueryRowx(insertQuery, map[string]interface{}{
			"Item":      insertInventory.Name,
			"NoOfItems": insertInventory.NoOfItems,
		}).StructScan(&insertedInventory)
	
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save inventory: " + err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"inventory": insertedInventory})
	})
	
}
