package model

type InquiryProductResp struct {
	ID             uint    `json:"id"`
	SKU            string  `json:"sku"`
	ProductName    string  `json:"product_name"`
	Price          float64 `json:"price"`
	PriceFormatted string  `json:"price_formatted"`
	Qty            int     `json:"qty"`
}
