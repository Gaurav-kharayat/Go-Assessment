package models

type Inventory struct {
	ID         int     `json:"id"`
	ItemName   string  `json:"item_name"`
	SKU        string  `json:"sku"`
	StockCount int     `json:"stock_count"`
	Price      float64 `json:"price"`
	UpdatedAt  string  `json:"updated_at"`
}
