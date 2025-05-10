package metalase

import (
	"base-golang/internal/entity"
	"base-golang/internal/model"
)

type MEtalaseMapper interface {
	MapAddEtalaseRequestToEntity(req model.AddEtalaseRequest) entity.Etalase
	ContainSpecialChars(str string) bool
}
