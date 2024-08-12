package internal

import (
	"context"
	"github.com/antlko/goauth-boilerplate/internal/db"
	"github.com/antlko/goauth-boilerplate/internal/logger"
	oauthgoogle "github.com/antlko/goauth-boilerplate/internal/oauth2/google"
	"github.com/antlko/goauth-boilerplate/internal/server"
	"log/slog"
)

type AppConfig struct {
	Hostname        string `env:"HOSTNAME"`
	ApplicationName string `env:"APPLICATION_NAME"`

	Server       server.Config
	DB           db.Config
	GoogleOauth2 oauthgoogle.Config
}

func InitService(cfg AppConfig) {
	ctx := context.Background()

	logger.InitLogger(logger.Config{
		AppName:  cfg.ApplicationName,
		Hostname: cfg.Hostname,
	})

	dbInst, err := db.NewDB(cfg.DB, cfg.ApplicationName)
	if err != nil {
		slog.ErrorContext(ctx, "db initialisation: %s", err.Error())
		return
	}

	googleConfig := oauthgoogle.InitConfig(cfg.GoogleOauth2)

	if err := server.InitServer(cfg.Server, dbInst, googleConfig); err != nil {
		slog.ErrorContext(ctx, "server initialisation: %s", err.Error())
	}
}
