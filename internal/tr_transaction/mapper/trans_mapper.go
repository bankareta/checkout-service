package mapper

import (
	"checkout-service/internal/entity"
	"checkout-service/internal/model"
	tr_transaction "checkout-service/internal/tr_transaction"
	"fmt"
	"math"
	"strings"

	"github.com/sirupsen/logrus"
)

type transactionMapper struct {
	Logger *logrus.Logger
}

func NewTransactionMapper(logger *logrus.Logger) tr_transaction.TransactionMapper {
	return &transactionMapper{Logger: logger}
}

func (m transactionMapper) MapScanProductsResp(resData model.ScanProductResponse, res entity.Products) model.ScanProductResponse {
	return model.ScanProductResponse{
		ID:             res.ID,
		SKU:            res.SKU,
		ProductName:    res.Name,
		Price:          res.Price,
		PriceFormatted: m.FormatPrice(res.Price),
		Qty:            res.InventoryQty,
	}
}

func (m transactionMapper) FormatPrice(price float64) string {
	s := fmt.Sprintf("%.2f", price)
	parts := strings.Split(s, ".")
	intPart := parts[0]
	decimalPart := parts[1]

	var result strings.Builder
	n := len(intPart)
	for i, c := range intPart {
		if (n-i)%3 == 0 && i != 0 {
			result.WriteRune(',')
		}
		result.WriteRune(c)
	}

	return "$" + result.String() + "." + decimalPart
}
func (m transactionMapper) RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
