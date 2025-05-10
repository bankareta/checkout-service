package metalase

import (
	"checkout-service/internal/model"
	"context"
)

type MEtalaseUseCase interface {
	AddEtalase(ctx context.Context, request *model.AddEtalaseRequest) error
}
