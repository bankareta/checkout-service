package master_products

import (
	"checkout-service/internal/entity"
	"checkout-service/internal/model"
)

type MasterProductsMapper interface {
	MapInquiryProductsResp(resData []model.InquiryProductResp, res []entity.Products) []model.InquiryProductResp
	FormatPrice(price float64) string
}
