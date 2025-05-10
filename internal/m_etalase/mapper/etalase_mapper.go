package mapper

import (
	"checkout-service/internal/entity"
	metalase "checkout-service/internal/m_etalase"
	"checkout-service/internal/model"

	"github.com/sirupsen/logrus"
)

type metalaseMapper struct {
	Logger *logrus.Logger
}

func NewMEtalaseMapper(logger *logrus.Logger) metalase.MEtalaseMapper {
	return &metalaseMapper{Logger: logger}
}

func (m metalaseMapper) MapAddEtalaseRequestToEntity(req model.AddEtalaseRequest) entity.Etalase {
	return entity.Etalase{
		EtalaseName:  req.EtalaseName,
		ReferenceId:  req.ReferenceId,
		Username:     req.Username,
		PhotoEtalase: req.PhotoEtalase,
		DeletedAt:    nil,
	}
}

func (m metalaseMapper) ContainSpecialChars(str string) bool {
	disallowedChars := []rune{'“', '$', '%', '‘', ';', '<', '>', '[', ']', '\\', '\''}

	for _, char := range str {
		for _, disallowed := range disallowedChars {
			if char == disallowed {
				return true
			}
		}
	}
	return false
}
