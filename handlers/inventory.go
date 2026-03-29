package handlers

import (
	"database/sql"
	"log"
	"time"

	"inventory-service/db"
	"inventory-service/models"

	"github.com/gin-gonic/gin"
)

type StockUpdateRequest struct {
	SKU        string `json:"sku"`
	Adjustment int    `json:"adjustment"`
}

func GetInventory(c *gin.Context) {
	rows, err := db.DB.Query("SELECT * FROM inventory")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var items []models.Inventory

	for rows.Next() {
		var item models.Inventory
		err := rows.Scan(&item.ID, &item.ItemName, &item.SKU, &item.StockCount, &item.Price, &item.UpdatedAt)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		items = append(items, item)
	}

	c.JSON(200, items)
}

func UpdateStock(c *gin.Context) {
	var req StockUpdateRequest

	// 🔹 Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	// 🔥 Validation
	if req.SKU == "" {
		c.JSON(400, gin.H{"error": "SKU is required"})
		return
	}

	if req.Adjustment == 0 {
		c.JSON(400, gin.H{"error": "Adjustment is required and cannot be 0"})
		return
	}

	// 🔹 Get current stock
	var currentStock int

	err := db.DB.QueryRow(
		"SELECT stock_count FROM inventory WHERE sku=$1",
		req.SKU,
	).Scan(&currentStock)

	if err == sql.ErrNoRows {
		c.JSON(404, gin.H{"error": "SKU not found"})
		return
	} else if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 🔹 Calculate new stock
	newStock := currentStock + req.Adjustment

	// 🔥 Prevent negative stock
	if newStock < 0 {
		c.JSON(400, gin.H{"error": "Stock cannot be negative"})
		return
	}

	// 🔹 Update DB
	_, err = db.DB.Exec(
		"UPDATE inventory SET stock_count=$1, updated_at=NOW() WHERE sku=$2",
		newStock, req.SKU,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 🔥 Logging (IMPORTANT for assignment)
	log.Printf("[LOG] SKU: %s | Old Stock: %d | New Stock: %d | Time: %s\n",
		req.SKU, currentStock, newStock, time.Now().Format(time.RFC3339))

	// 🔹 Success response
	c.JSON(200, gin.H{
		"message":   "Stock updated",
		"new_stock": newStock,
	})
}

func GetLowStock(c *gin.Context) {
	rows, err := db.DB.Query("SELECT * FROM inventory WHERE stock_count < 10")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var items []models.Inventory

	for rows.Next() {
		var item models.Inventory
		err := rows.Scan(
			&item.ID,
			&item.ItemName,
			&item.SKU,
			&item.StockCount,
			&item.Price,
			&item.UpdatedAt,
		)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		items = append(items, item)
	}

	c.JSON(200, items)
}
