package route

import (
	"checkout-service/internal/delivery/http/middleware"
	metalaseHttp "checkout-service/internal/m_etalase/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App               *fiber.App
	EtalaseController *metalaseHttp.MEtalaseController
}

func (c *RouteConfig) Setup() {
	c.App.Use(middleware.AddHeader())
	c.App.Use(middleware.CustomLogger())
	c.SetupGuestRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	etalase := c.App.Group("/api/v1/etalase")
	etalase.Post("/addEtalase", c.EtalaseController.AddEtalase)
}
