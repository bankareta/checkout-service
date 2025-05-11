package http

import (
	master_products "checkout-service/internal/master_products"
	"checkout-service/internal/model"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type MasterProductsController struct {
	Log     *logrus.Logger
	UseCase master_products.MasterProductsUseCase
}

func NewMasterProductsController(useCase master_products.MasterProductsUseCase, logger *logrus.Logger) *MasterProductsController {
	return &MasterProductsController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *MasterProductsController) InquiryProducts(ctx *fiber.Ctx) error {
	// request := new(model.InquiryProductsRequest)
	// err := ctx.BodyParser(request)
	// if err != nil {
	// 	c.Log.Warnf("Failed to parse request body : %+v", err)
	// 	return fiber.ErrBadRequest
	// }

	res, err := c.UseCase.InquiryProducts()
	if err != nil {
		c.Log.Warnf("Failed to get products : %+v", err)
		return err
	}
	return model.ResponseSuccess[any](ctx, "Success Get Products", res, nil)

}
