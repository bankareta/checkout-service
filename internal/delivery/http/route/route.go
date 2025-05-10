package route

import (
	"base-golang/internal/delivery/http/middleware"
	metalaseHttp "base-golang/internal/m_etalase/delivery/http"

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
