package model

import "time"

type ScanProductRequest struct {
	SKU string `json:"sku"`
}
type ScanProductResponse struct {
	ID             uint    `json:"id"`
	SKU            string  `json:"sku"`
	ProductName    string  `json:"product_name"`
	Price          float64 `json:"price"`
	PriceFormatted string  `json:"price_formatted"`
	Qty            int     `json:"qty"`
}

type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}

type CheckoutItem struct {
	ProductID int `json:"product_id"`
	Qty       int `json:"qty"`
}

type CheckoutResponse struct {
	ID             int                    `json:"id"`
	CheckoutDate   time.Time              `json:"checkout_date"`
	Items          []CheckoutItemResponse `json:"items"`
	PriceTotal     float64                `json:"price_total"`
	TotalFormatted string                 `json:"total_formatted"`
}

type CheckoutItemResponse struct {
	Products            []ProductCheckoutResponse `json:"products"`
	PriceTotal          float64                   `json:"price_total"`
	TotalFormatted      string                    `json:"total_formatted"`
	DiscountDescription []string                  `json:"discount_description"`
}

type ProductCheckoutResponse struct {
	ProductID      int     `json:"product_id"`
	ProductName    string  `json:"product_name"`
	Qty            int     `json:"qty"`
	Price          float64 `json:"price"`
	PriceTotal     float64 `json:"price_total"`
	PriceFormatted string  `json:"price_formatted"`
	TotalFormatted string  `json:"total_formatted"`
}
