package master_products

import (
	"checkout-service/internal/model"
)

type MasterProductsUseCase interface {
	InquiryProducts() ([]model.InquiryProductResp, error)
}
