package entity

import "time"

// User is a struct that represents a user entity
type Transaction struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	CustomerName  *string   `gorm:"column:customer_name" json:"customer_name"`
	CustomerPhone *string   `gorm:"column:customer_phone" json:"customer_phone"`
	Status        int       `gorm:"column:status" json:"status"`
	TotalPrice    float64   `gorm:"column:total_price" json:"total_price"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"-"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"-"`
}

type TransactionDetail struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	TransactionID int       `gorm:"column:transaction_id" json:"transaction_id"`
	ProductID     int       `gorm:"column:product_id" json:"product_id"`
	ProductName   *string   `gorm:"column:product_name" json:"product_name"`
	SKU           *string   `gorm:"column:sku" json:"sku"`
	Qty           int       `gorm:"column:qty" json:"qty"`
	Price         float64   `gorm:"column:price" json:"price"`
	Status        int       `gorm:"column:status" json:"status"`
	TotalPrice    float64   `gorm:"column:total_price" json:"total_price"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"-"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"-"`
}

type TransactionDetailDiscount struct {
	ID                      uint      `gorm:"primaryKey" json:"id"`
	DetailTransactionID     int       `gorm:"column:detail_transaction_id" json:"detail_transaction_id"`
	DiscountID              int       `gorm:"column:discount_id" json:"discount_id"`
	Type                    int       `gorm:"column:type" json:"type"`
	IsPercentage            int       `gorm:"column:is_percentage" json:"is_percentage"`
	Amount                  int       `gorm:"column:amount" json:"amount"`
	RequiredQty             int       `gorm:"column:required_qty" json:"required_qty"`
	FinalQty                int       `gorm:"column:final_qty" json:"final_qty"`
	DetailTransactionFreeID int       `gorm:"column:detail_transaction_free_id" json:"detail_transaction_free_id"`
	ProductID               int       `gorm:"column:product_id" json:"product_id"`
	ProductName             string    `gorm:"column:product_name" json:"product_name"`
	DiscountDescription     string    `gorm:"column:discount_description" json:"discount_description"`
	CreatedAt               time.Time `gorm:"column:created_at" json:"-"`
	UpdatedAt               time.Time `gorm:"column:updated_at" json:"-"`
}

func (u *Transaction) TableName() string {
	return "transaction"
}

func (u *TransactionDetail) TableName() string {
	return "transaction_detail"
}

func (u *TransactionDetailDiscount) TableName() string {
	return "transaction_detail_discount"
}
