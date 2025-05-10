package metalase

import (
	"checkout-service/internal/entity"
)

type MEtalaseRepository interface {
	PostEtalase(dataAddEtalase entity.Etalase) (entity.Etalase, error)
}
