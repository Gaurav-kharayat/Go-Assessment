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

// GetInventory retrieves all inventory items from the database
func GetInventory(c *gin.Context) {

	// Execute query to fetch all records from inventory table
	rows, err := db.DB.Query("SELECT * FROM inventory")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var items []models.Inventory

	// Iterate through all rows returned by query
	for rows.Next() {
		var item models.Inventory
		err := rows.Scan(&item.ID, &item.ItemName, &item.SKU, &item.StockCount, &item.Price, &item.UpdatedAt)
		if err != nil {
			// Handle scanning errors
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		items = append(items, item)
	}

	c.JSON(200, items)
}

func UpdateStock(c *gin.Context) {
	var req StockUpdateRequest

	// Bind incoming JSON payload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	// Validate required fields
	if req.SKU == "" {
		c.JSON(400, gin.H{"error": "SKU is required"})
		return
	}

	if req.Adjustment == 0 {
		c.JSON(400, gin.H{"error": "Adjustment is required and cannot be 0"})
		return
	}

	// Check if SKU exists first so we can return proper 404
	var exists bool
	err := db.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM inventory WHERE sku = $1)",
		req.SKU,
	).Scan(&exists)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if !exists {
		c.JSON(404, gin.H{"error": "SKU not found"})
		return
	}

	// Atomic update to avoid race conditions during concurrent requests
	query := `
	UPDATE inventory
	SET stock_count = stock_count + $1,
	    updated_at = NOW()
	WHERE sku = $2 AND stock_count + $1 >= 0
	RETURNING stock_count;
	`

	var newStock int

	err = db.DB.QueryRow(query, req.Adjustment, req.SKU).Scan(&newStock)

	// If no row returned here, SKU exists but stock would go below zero
	if err == sql.ErrNoRows {
		c.JSON(400, gin.H{"error": "Stock cannot go below zero"})
		return
	} else if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Log successful stock update
	log.Printf("[LOG] SKU: %s | Adjustment: %d | New Stock: %d | Time: %s\n",
		req.SKU, req.Adjustment, newStock, time.Now().Format(time.RFC3339))

	c.JSON(200, gin.H{
		"message":   "Stock updated",
		"new_stock": newStock,
	})
}

// GetLowStock retrieves all inventory items where stock_count is less than 10
func GetLowStock(c *gin.Context) {

	// Query database for items with low stock
	rows, err := db.DB.Query("SELECT * FROM inventory WHERE stock_count < 10")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var items []models.Inventory

	// Iterate through query results
	for rows.Next() {
		var item models.Inventory
		// Scan each row into Inventory struct
		err := rows.Scan(
			&item.ID,
			&item.ItemName,
			&item.SKU,
			&item.StockCount,
			&item.Price,
			&item.UpdatedAt,
		)
		if err != nil {
			// Handle scanning errors
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		// Append item to response list
		items = append(items, item)
	}

	c.JSON(200, items)
}
