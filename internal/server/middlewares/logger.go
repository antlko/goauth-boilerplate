package middlewares

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"log/slog"
)

func Logger(c fiber.Ctx) error {
	if err := c.Next(); err != nil {
		slog.ErrorContext(c.Context(), fmt.Sprintf("server request error: %s", err.Error()))
	}

	slog.InfoContext(c.Context(),
		"new http request",
		"ip", c.IP(),
		"url", c.OriginalURL(),
		"request_body", string(c.Body()),
		"response_body", string(c.Response().Body()),
	)
	return nil
}
