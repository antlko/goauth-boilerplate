package handlers

import (
	"context"
	"github.com/antlko/goauth-boilerplate/internal/db"
	"github.com/antlko/goauth-boilerplate/internal/server/responses"
	"github.com/gofiber/fiber/v3"
	"log/slog"
	"net/http"
)

type userGetterByLogin interface {
	GetByLogin(ctx context.Context, login string) (db.User, error)
}

type UserHandler struct {
	userGetterByLogin userGetterByLogin
}

func NewUserHandler(userGetterByLogin userGetterByLogin) UserHandler {
	return UserHandler{
		userGetterByLogin: userGetterByLogin,
	}
}

func (h UserHandler) GetUser(c fiber.Ctx) error {
	ctx := c.Context()
	login := c.GetReqHeaders()["X-User-Id"][0]
	user, err := h.userGetterByLogin.GetByLogin(ctx, login)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return c.Status(http.StatusInternalServerError).JSON(responses.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "can't make a fetch",
		})
	}
	return c.Status(http.StatusOK).JSON(responses.User{
		Id:    user.Id,
		Login: user.Login,
		Email: user.Email,
	})
}
