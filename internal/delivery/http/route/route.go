package route

import (
	"checkout-service/internal/delivery/http/middleware"
	masterProductsHttp "checkout-service/internal/master_products/delivery/http"
	transactionHttp "checkout-service/internal/tr_transaction/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                      *fiber.App
	MasterProductsController *masterProductsHttp.MasterProductsController
	TransactionController    *transactionHttp.TransactionController
}

func (c *RouteConfig) Setup() {
	c.App.Use(middleware.AddHeader())
	c.App.Use(middleware.CustomLogger())
	c.SetupGuestRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	master := c.App.Group("/api/v1/master")
	master.Post("/inquiryProducts", c.MasterProductsController.InquiryProducts)

	trans := c.App.Group("/api/v1/cashier")
	trans.Post("/scanProduct", c.TransactionController.ScanProduct)
	trans.Post("/checkout", c.TransactionController.Checkout)
}
