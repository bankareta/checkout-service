package http

import (
	metalase "base-golang/internal/m_etalase"
	"base-golang/internal/model"

	"github.com/gofiber/fiber/v2"

	"github.com/sirupsen/logrus"
)

type MEtalaseController struct {
	Log     *logrus.Logger
	UseCase metalase.MEtalaseUseCase
}

func NewMEtalaseController(useCase metalase.MEtalaseUseCase, logger *logrus.Logger) *MEtalaseController {
	return &MEtalaseController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *MEtalaseController) AddEtalase(ctx *fiber.Ctx) error {
	request := new(model.AddEtalaseRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}
	err = c.UseCase.AddEtalase(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to Add Etalase: %+v", err)
		return err
	}
	return model.ResponseSuccess[any](ctx, "Success Add Etalase", "Success", nil)
}
