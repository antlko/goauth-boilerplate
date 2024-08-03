package main

import (
	"context"
	"github.com/antlko/goauth-boilerplate/internal"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"log/slog"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error(".env not load")
	}

	var cfg internal.AppConfig
	ctx := context.Background()

	if err := envconfig.Process(ctx, &cfg); err != nil {
		slog.Error("process config: %s", err)
		return
	}
	internal.InitService(cfg)
}
