package config

import (
	"checkout-service/internal/environtment"
	"checkout-service/internal/helpers"
	"checkout-service/internal/model"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func NewFiber(config *environtment.Config) *fiber.App {
	var app = fiber.New(fiber.Config{
		AppName:      config.APP_NAME,
		ErrorHandler: NewErrorHandler(),
		Prefork:      config.WEB_PREFORK,
	})

	app.Use(recover.New())

	return app
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		rc := "99"
		message := "General Error"

		fiturCheck := ""
		fitur := ctx.UserContext().Value("fitur")
		var errData any
		if e, ok := err.(*helpers.Error); ok {
			code = e.Code
			rc = e.Rc
			message = e.Message
			errData = e.Errors
		}
		if fitur != nil {
			fiturCheck = fitur.(string)
		}
		fmt.Println(fiturCheck)
		return model.ResponseError[any](ctx, code, rc, message, errData)
	}
}
