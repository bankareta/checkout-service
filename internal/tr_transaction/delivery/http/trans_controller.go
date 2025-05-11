package http

import (
	"checkout-service/internal/model"
	tr_transaction "checkout-service/internal/tr_transaction"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type TransactionController struct {
	Log     *logrus.Logger
	UseCase tr_transaction.TransactionUseCase
}

func NewTransactionController(useCase tr_transaction.TransactionUseCase, logger *logrus.Logger) *TransactionController {
	return &TransactionController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *TransactionController) ScanProduct(ctx *fiber.Ctx) error {
	request := new(model.ScanProductRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	res, err := c.UseCase.ScanProduct(ctx.UserContext(), *request)
	if err != nil {
		c.Log.Warnf("Failed to get products : %+v", err)
		return err
	}
	return model.ResponseSuccess[any](ctx, "Success Get Products", res, nil)
}

func (c *TransactionController) Checkout(ctx *fiber.Ctx) error {
	request := new(model.CheckoutRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	res, err := c.UseCase.Checkout(ctx.UserContext(), *request)
	if err != nil {
		c.Log.Warnf("Failed to get products : %+v", err)
		return err
	}
	return model.ResponseSuccess[any](ctx, "Success Checkout", res, nil)
}
