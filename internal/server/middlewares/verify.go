package middlewares

import (
	"github.com/antlko/goauth-boilerplate/internal/server/responses"
	"github.com/gofiber/fiber/v3"
	"net/http"
)

type tokenValidator interface {
	Validate(token string) (bool, string, error)
}

func BearerVerifier(tokenValidator tokenValidator) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		authHeaders, ok := c.GetReqHeaders()["Authorization"]
		if !ok {
			return c.Status(http.StatusUnauthorized).JSON(responses.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "bad auth header",
			})
		}
		if len(authHeaders) == 0 {
			return c.Status(http.StatusUnauthorized).JSON(responses.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "bad auth header",
			})
		}
		authHeader := authHeaders[0]
		if len(authHeader) < 10 {
			return c.Status(http.StatusUnauthorized).JSON(responses.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "bad auth header",
			})
		}
		authHeader = authHeader[len("Bearer "):]

		isValid, username, err := tokenValidator.Validate(authHeader)
		if err != nil || !isValid {
			return c.Status(http.StatusUnauthorized).JSON(responses.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "token not valid",
			})
		}

		c.Request().Header.Add("X-User-Id", username)
		return c.Next()
	}
}
