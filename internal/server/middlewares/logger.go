package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"log/slog"
)

const secretDataTemplate = "**********"

func Logger(c fiber.Ctx) error {
	ctx := c.Context()

	if err := c.Next(); err != nil {
		slog.ErrorContext(c.Context(), fmt.Sprintf("server request error: %s", err.Error()))
	}

	body := hidePII(ctx, c.Body())

	slog.InfoContext(ctx,
		"new http request",
		"ip", c.IP(),
		"url", c.OriginalURL(),
		"request_body", body,
		"response_body", string(c.Response().Body()),
	)
	return nil
}

// Simple example to hide from the request logs user password
func hidePII(ctx context.Context, data []byte) string {
	var body string
	if data == nil || len(data) == 0 {
		return body
	}
	result := make(map[string]any)
	if err := json.Unmarshal(data, &result); err != nil {
		slog.ErrorContext(ctx, "unmarshalling request body", err.Error())
		return body
	}
	if len(result) > 0 {
		if _, ok := result["password"]; ok {
			result["password"] = secretDataTemplate
		}
		bodyData, _ := json.Marshal(result)
		body = string(bodyData)
	}
	return body
}
