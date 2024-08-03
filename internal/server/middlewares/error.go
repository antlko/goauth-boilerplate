package middlewares

import (
	"fmt"
	"github.com/antlko/goauth-boilerplate/internal/server/responses"
	"github.com/gofiber/fiber/v3"
	"log/slog"
	"net/http"
)

func Error(c fiber.Ctx) error {
	defer func() {
		if err := recover(); err != nil {
			slog.ErrorContext(c.Context(), fmt.Sprintf("server handled panic error: %s", err))
			_ = c.Status(http.StatusInternalServerError).JSON(responses.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			})
		}
	}()

	if err := c.Next(); err != nil {
		slog.ErrorContext(c.Context(), fmt.Sprintf("server requested error: %s", err.Error()))
		return c.Status(http.StatusInternalServerError).JSON(responses.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}
	return nil
}
