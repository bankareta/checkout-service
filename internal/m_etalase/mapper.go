package metalase

import (
	"checkout-service/internal/entity"
	"checkout-service/internal/model"
)

type MEtalaseMapper interface {
	MapAddEtalaseRequestToEntity(req model.AddEtalaseRequest) entity.Etalase
	ContainSpecialChars(str string) bool
}
