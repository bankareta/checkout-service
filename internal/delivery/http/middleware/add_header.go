package middleware

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AddHeader() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		uuidString := uuid.New().String()
		ctx.Set("Request-ID", uuidString)
		ctx.Request().Header.Set("Request-ID", uuidString)
		newCtx := context.WithValue(ctx.UserContext(), "requestHeader", ctx.GetReqHeaders())
		ctx.SetUserContext(newCtx)
		return ctx.Next()
	}
}
