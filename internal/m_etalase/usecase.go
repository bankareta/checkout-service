package metalase

import (
	"base-golang/internal/model"
	"context"
)

type MEtalaseUseCase interface {
	AddEtalase(ctx context.Context, request *model.AddEtalaseRequest) error
}
