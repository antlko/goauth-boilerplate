package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/antlko/goauth-boilerplate/internal/db"
	"github.com/antlko/goauth-boilerplate/internal/jwt"
	"github.com/antlko/goauth-boilerplate/internal/server/requests"
	"github.com/antlko/goauth-boilerplate/internal/server/responses"
	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
)

type (
	userInserter interface {
		Insert(ctx context.Context, user db.User) error
	}
	userGetter interface {
		GetByLoginOrEmail(ctx context.Context, login, email string) (db.User, error)
		GetByLogin(ctx context.Context, login string) (db.User, error)
	}
	authorizer interface {
		CreateTokens(username string) (jwt.Tokens, error)
		ValidateAndUpdate(refresh string) (jwt.Tokens, error)
	}
)

type AuthHandler struct {
	userInserter userInserter
	userGetter   userGetter
	authorizer   authorizer
}

func NewAuthHandler(userInserter userInserter, userGetter userGetter, authorizer authorizer) AuthHandler {
	return AuthHandler{
		userInserter: userInserter,
		userGetter:   userGetter,
		authorizer:   authorizer,
	}
}

func (a AuthHandler) SignUp(c fiber.Ctx) error {
	ctx := c.Context()

	var request requests.SignUpRequest
	if err := json.Unmarshal(c.Body(), &request); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return c.Status(http.StatusBadRequest).JSON(responses.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "request body not parsed",
		})
	}

	user, err := a.userGetter.GetByLoginOrEmail(ctx, request.Login, request.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		slog.ErrorContext(ctx, err.Error())
		return c.Status(http.StatusInternalServerError).JSON(responses.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "can't make a fetch",
		})
	}
	if user.Email == request.Email || user.Login == request.Login {
		return c.Status(http.StatusBadRequest).JSON(responses.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "user with this email or login already exists",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 8)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return c.Status(http.StatusInternalServerError).JSON(responses.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	if err := a.userInserter.Insert(ctx, db.User{
		Login:    request.Login,
		Email:    request.Email,
		Password: string(hashedPassword),
	}); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return c.Status(http.StatusInternalServerError).JSON(responses.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "user not saved",
		})
	}

	return c.Status(http.StatusOK).JSON(responses.StatusResponse{
		Status: "ok",
	})
}

func (a AuthHandler) SignIn(c fiber.Ctx) error {
	ctx := c.Context()

	var request requests.SignInRequest
	if err := json.Unmarshal(c.Body(), &request); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return c.Status(http.StatusBadRequest).JSON(responses.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "request body not parsed",
		})
	}

	user, err := a.userGetter.GetByLogin(ctx, request.Login)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		slog.ErrorContext(ctx, err.Error())
		return c.Status(http.StatusInternalServerError).JSON(responses.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "can't make a fetch",
		})
	}
	if errors.Is(err, sql.ErrNoRows) {
		slog.ErrorContext(ctx, err.Error())
		return c.Status(http.StatusBadRequest).JSON(responses.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "incorrect login or password",
		})
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return c.Status(http.StatusBadRequest).JSON(responses.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "incorrect login or password",
		})
	}

	tokens, err := a.authorizer.CreateTokens(user.Login)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return c.Status(http.StatusInternalServerError).JSON(responses.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	return c.Status(http.StatusOK).JSON(responses.TokensResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

func (a AuthHandler) Verify(c fiber.Ctx) error {
	ctx := c.Context()

	var request requests.VerifyAndRefreshRequest
	if err := json.Unmarshal(c.Body(), &request); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return c.Status(http.StatusBadRequest).JSON(responses.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "request body not parsed",
		})
	}

	tokens, err := a.authorizer.ValidateAndUpdate(request.RefreshToken)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return c.Status(http.StatusUnauthorized).JSON(responses.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		})
	}
	return c.Status(http.StatusOK).JSON(responses.TokensResponse{
		AccessToken: tokens.AccessToken,
	})
}
