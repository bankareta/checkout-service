package entity

import "time"

// User is a struct that represents a user entity
type Discount struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	Type          int        `gorm:"column:type" json:"type"`
	IsPercentage  int        `gorm:"column:is_percentage" json:"is_percentage"`
	Amount        int        `gorm:"column:amount" json:"amount"`
	RequiredQty   int        `gorm:"column:required_qty" json:"required_qty"`
	FinalQty      int        `gorm:"column:final_qty" json:"final_qty"`
	FreeIDProduct int        `gorm:"column:free_id_product" json:"free_id_product"`
	CreatedAt     time.Time  `gorm:"column:created_at" json:"-"`
	UpdatedAt     time.Time  `gorm:"column:updated_at" json:"-"`
	DeletedAt     *time.Time `gorm:"column:deleted_at" json:"-"`
}

type ProductsDiscount struct {
	ProductId   int `gorm:"column:product_id" json:"product_id"`
	DiscountsId int `gorm:"column:discounts_id" json:"discounts_id"`
}

func (u *Discount) TableName() string {
	return "discounts"
}

func (u *ProductsDiscount) TableName() string {
	return "products_discounts"
}
