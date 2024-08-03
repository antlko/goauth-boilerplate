package db

import (
	"embed"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type Config struct {
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	DBName   string `env:"DB_NAME"`
	Schema   string `env:"DB_SCHEMA"`

	MaxIdleConns int    `env:"DB_MAX_IDLE_CONNS" envDefault:"2"`
	MaxOpenConns int    `env:"DB_MAX_OPEN_CONNS" envDefault:"4"`
	LogLevel     string `env:"DB_LOG_LEVEL" envDefault:"error"`
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func NewDB(cfg Config, appName string) (*sqlx.DB, error) {
	db, err := sqlx.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable application_name=%s search_path=%s",
			cfg.Host,
			cfg.Port,
			cfg.User,
			cfg.Password,
			cfg.DBName,
			appName,
			cfg.Schema,
		))
	if err != nil {
		return nil, fmt.Errorf("db conn error: %w", err)
	}

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	if _, err = db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", cfg.Schema)); err != nil {
		return nil, fmt.Errorf("create schema: %w", err)
	}

	goose.SetBaseFS(embedMigrations)

	if err = goose.SetDialect("postgres"); err != nil {
		return nil, fmt.Errorf("set migrations dialect: %w", err)
	}

	if err = goose.Up(db.Unsafe().DB, "migrations"); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	return db.Unsafe(), nil
}
