package server

import (
	"fmt"
	"github.com/antlko/goauth-boilerplate/internal/db"
	"github.com/antlko/goauth-boilerplate/internal/jwt"
	"github.com/antlko/goauth-boilerplate/internal/server/handlers"
	"github.com/antlko/goauth-boilerplate/internal/server/middlewares"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	ServerPort string `env:"SERVER_PORT"`
	JwtConfig  jwt.Config
}

func InitServer(cfg Config, dbInst *sqlx.DB) error {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Accept-Language", "Content-Length", "Authorization"},
		AllowOrigins: []string{"*", "Access-Control-Allow-Headers"},
	}))

	userRepo := db.NewUserRepo(dbInst)
	authorizer := jwt.NewAuthorizer(cfg.JwtConfig)

	authHandler := handlers.NewAuthHandler(userRepo, userRepo, authorizer)
	userHandler := handlers.NewUserHandler(userRepo)

	app.Use(
		middlewares.Logger,
		middlewares.Error,
	)

	app.Post("/api/v1/users/auth/signup", authHandler.SignUp)
	app.Post("/api/v1/users/auth/signin", authHandler.SignIn)
	app.Post("/api/v1/users/auth/token/refresh", authHandler.Verify)

	protected := app.Group("/api/v1/protected", middlewares.BearerVerifier(authorizer))
	protected.Get("/user", userHandler.GetUser)

	if err := app.Listen(":" + cfg.ServerPort); err != nil {
		return fmt.Errorf("server listen: %w", err)
	}
	return nil
}
