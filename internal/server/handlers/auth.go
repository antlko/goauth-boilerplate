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
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"log/slog"
	"net/http"
	"time"
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
	googleAuthorizer interface {
		Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
		AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
		Client(ctx context.Context, t *oauth2.Token) *http.Client
	}
)

type AuthHandler struct {
	userInserter     userInserter
	userGetter       userGetter
	authorizer       authorizer
	googleAuthorizer googleAuthorizer

	sessionStore *session.Store
}

func NewAuthHandler(userInserter userInserter, userGetter userGetter, authorizer authorizer, googleConfig googleAuthorizer) AuthHandler {
	return AuthHandler{
		userInserter:     userInserter,
		userGetter:       userGetter,
		authorizer:       authorizer,
		googleAuthorizer: googleConfig,
		sessionStore:     session.New(),
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

func (a AuthHandler) GoogleSignIn(c fiber.Ctx) error {
	stateId := uuid.NewString()
	ctx := c.Context()

	sess, err := a.sessionStore.Get(c)
	if err != nil {
		slog.ErrorContext(ctx, "session not initialized")
		return c.Status(http.StatusInternalServerError).JSON(responses.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "session not initialized",
		})
	}

	sess.Set("state", stateId)
	sess.SetExpiry(time.Minute * 5)
	if err := sess.Save(); err != nil {
		slog.ErrorContext(ctx, "saving session")
		return c.Status(http.StatusInternalServerError).JSON(responses.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "session storing error",
		})
	}

	url := a.googleAuthorizer.AuthCodeURL(stateId)
	return c.Status(http.StatusTemporaryRedirect).Redirect().To(url)
}

func (a AuthHandler) GoogleCallback(c fiber.Ctx) error {
	ctx := c.Context()

	sess, err := a.sessionStore.Get(c)
	if err != nil {
		slog.ErrorContext(ctx, "invalid oauth state")
		return c.Status(http.StatusInternalServerError).JSON(responses.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "session not initialized",
		})
	}

	stateId := sess.Get("state")
	if c.FormValue("state") != stateId {
		slog.ErrorContext(ctx, "invalid oauth state")
		return c.Status(http.StatusUnauthorized).JSON(responses.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		})
	}

	code := c.FormValue("code")
	token, err := a.googleAuthorizer.Exchange(ctx, code)
	if err != nil {
		slog.ErrorContext(ctx, "exchange code", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(responses.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "failed to get the token",
		})
	}

	client := a.googleAuthorizer.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		slog.ErrorContext(ctx, "get userinfo", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(responses.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "failed to get google client",
		})
	}
	defer resp.Body.Close()

	var userInfo responses.GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		slog.ErrorContext(ctx, "parse userinfo", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(responses.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "failed to extract the user",
		})
	}

	_, err = a.userGetter.GetByLogin(ctx, userInfo.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		slog.ErrorContext(ctx, err.Error())
		return c.Status(http.StatusInternalServerError).JSON(responses.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "can't make a fetch",
		})
	}
	if errors.Is(err, sql.ErrNoRows) {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(uuid.NewString()), 8)
		if err != nil {
			slog.ErrorContext(ctx, err.Error())
			return c.Status(http.StatusInternalServerError).JSON(responses.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			})
		}

		if err := a.userInserter.Insert(ctx, db.User{
			Login:    uuid.NewString(),
			Email:    userInfo.Email,
			Password: string(hashedPassword),
		}); err != nil {
			slog.ErrorContext(ctx, err.Error())
			return c.Status(http.StatusInternalServerError).JSON(responses.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "user not saved",
			})
		}
	}

	tokens, err := a.authorizer.CreateTokens(userInfo.Email)
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
