package v1

import (
	"github.com/Taya-san/portofolio-backend/internal/controller/v1/response"
	"github.com/gofiber/fiber/v2"
)

func errorResponse(ctx *fiber.Ctx, code int, msg string) error {
	return ctx.Status(code).JSON(response.Error{Error: msg})
}
