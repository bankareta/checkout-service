package tr_transaction

import (
	"checkout-service/internal/model"
	"context"
)

type TransactionUseCase interface {
	ScanProduct(ctx context.Context, request model.ScanProductRequest) (model.ScanProductResponse, error)
	Checkout(ctx context.Context, request model.CheckoutRequest) (model.CheckoutResponse, error)
}
