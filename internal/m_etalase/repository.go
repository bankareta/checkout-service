package metalase

import (
	"base-golang/internal/entity"
)

type MEtalaseRepository interface {
	PostEtalase(dataAddEtalase entity.Etalase) (entity.Etalase, error)
}
