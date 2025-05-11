package tr_transaction

import (
	"checkout-service/internal/entity"
	"checkout-service/internal/model"
)

type TransactionMapper interface {
	MapScanProductsResp(resData model.ScanProductResponse, res entity.Products) model.ScanProductResponse
	FormatPrice(price float64) string
	RoundFloat(val float64, precision uint) float64
}
