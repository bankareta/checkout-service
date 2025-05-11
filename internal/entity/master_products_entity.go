package entity

import "time"

// User is a struct that represents a user entity
type Products struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	SKU          string     `gorm:"column:sku" json:"sku"`
	Name         string     `gorm:"column:name" json:"name"`
	Price        float64    `gorm:"column:price" json:"price"`
	InventoryQty int        `gorm:"column:inventory_qty" json:"inventory_qty"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"-"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"-"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"-"`
}

func (u *Products) TableName() string {
	return "products"
}
