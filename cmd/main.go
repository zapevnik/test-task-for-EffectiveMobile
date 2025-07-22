package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"subscription-service/internal/config"
	httpDelivery "subscription-service/internal/delivery/http"
	"subscription-service/internal/logger"
	"subscription-service/internal/storage/postgres"
	"subscription-service/internal/usecase/subscription"
)


func main() {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "config/config.yaml"
	}

	cfg, err := config.LoadConfig(path)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	var level slog.Level

	switch cfg.LogLevel {
	case "debug":
			level = slog.LevelDebug
	case "info":
			level = slog.LevelInfo
	case "warn":
			level = slog.LevelWarn
	case "error":
			level = slog.LevelError
	default:
			level = slog.LevelInfo
	}

	logger.Init(level)
	slog.SetDefault(logger.Log)
	slog.Info("config loaded:", "path", path)

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
    slog.Error("failed to connect database", "error", err)
    os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
    slog.Error("failed to ping database", "error", err)
    os.Exit(1)
	}

	storage := postgres.NewSubscriptionStorage(db, logger.Log)
	service := subscription.NewService(storage, logger.Log)
	handler := httpDelivery.NewHandler(service, logger.Log)
	router := httpDelivery.NewRouter(handler, logger.Log)

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	slog.Info("server starting", "addr", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
    slog.Error("server failed", "error", err)
    os.Exit(1)
	}
}
