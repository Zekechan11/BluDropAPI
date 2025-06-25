package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type InventoryItem struct {
	Inventory_id	int    `json:"inventory_id"`
	Item      		string `json:"item"`
	No_of_items 	int    `json:"no_of_items"`
}

func InventoryRoutes(r *gin.Engine, db *sqlx.DB) {
	r.GET("/api/get_inventory", func(ctx *gin.Context) {
		var inventory []InventoryItem
		
		err := db.Select(&inventory, "SELECT * FROM inventory")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, inventory)
	})

	r.POST("/api/save_inventory", func(ctx *gin.Context) {
		var insertInventory struct {
			Item 		string `json:"item"`
			No_of_items string  `json:"no_of_items"`
		}

		if err := ctx.ShouldBindJSON(&insertInventory); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		insertQuery := `
    		INSERT INTO inventory (item, no_of_items) 
    		VALUES (:item, :no_of_items)`

		_, err := db.NamedExec(insertQuery, insertInventory)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save inventory: " + err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"Inventory": insertInventory})
	})
	

	r.PUT("/api/update_inventory/:id", func(ctx *gin.Context) {
		var updateInventory struct {
			Inventory_id	int    `json:"inventory_id"`
			Item      		string `json:"item"`
			No_of_items 	string `json:"no_of_items"`
		}
	
		if err := ctx.ShouldBindJSON(&updateInventory); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id := ctx.Param("id")
    	updateInventory.Inventory_id, _ = strconv.Atoi(id)
	
		updateQuery := `
			UPDATE inventory 
			SET item = :item, no_of_items = :no_of_items
			WHERE Inventory_id = :inventory_id`
	
		_, err := db.NamedExec(updateQuery, updateInventory)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update inventory: " + err.Error()})
			return
		}
	
		ctx.JSON(http.StatusOK, gin.H{"message": "Inventory updated successfully", "inventory": updateInventory})
	})
	

	r.DELETE("/api/delete_inventory/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
	
		_, err := db.Exec("DELETE FROM inventory WHERE inventory_id = ?", id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete inventory: " + err.Error()})
			return
		}
	
		ctx.JSON(http.StatusOK, gin.H{"message": "Inventory deleted successfully"})
	})
	
}
