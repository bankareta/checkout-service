package mapper

import (
	"checkout-service/internal/entity"
	master_products "checkout-service/internal/master_products"
	"checkout-service/internal/model"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type masterProductsMapper struct {
	Logger *logrus.Logger
}

func NewMasterProductsMapper(logger *logrus.Logger) master_products.MasterProductsMapper {
	return &masterProductsMapper{Logger: logger}
}

func (m masterProductsMapper) MapInquiryProductsResp(resData []model.InquiryProductResp, res []entity.Products) []model.InquiryProductResp {
	convertModel := make([]model.InquiryProductResp, len(res))
	for i, mProduct := range res {
		convertModel[i] = model.InquiryProductResp{
			ID:             mProduct.ID,
			SKU:            mProduct.SKU,
			ProductName:    mProduct.Name,
			Price:          mProduct.Price,
			PriceFormatted: m.FormatPrice(mProduct.Price),
			Qty:            mProduct.InventoryQty,
		}
	}
	return convertModel
}

func (m masterProductsMapper) FormatPrice(price float64) string {
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
