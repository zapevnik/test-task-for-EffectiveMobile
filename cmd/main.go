// @title Subscription Service API
// @version 1.0
// @description API для управления подписками пользователей
// @host localhost:8080
// @BasePath /
package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	_ "subscription-service/docs"
	"subscription-service/internal/config"
	httpDelivery "subscription-service/internal/delivery/http"
	"subscription-service/pkg/logger"

	"subscription-service/pkg/storage"
	"subscription-service/internal/storage/postgres"
	"subscription-service/internal/usecase/subscription"
)

const defaultConfigPath string = "config/config.yaml"

func main() {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = defaultConfigPath
	}

	cfg, err := config.LoadConfig(path)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}


	logger.Init(cfg.LogLevel)
	slog.SetDefault(logger.Log)
	slog.Info("config loaded:", "path", path)

	db := storage.NewPostgresDB(storage.Config{
    Host:     cfg.Database.Host,
    Port:     cfg.Database.Port,
    User:     cfg.Database.User,
    Password: cfg.Database.Password,
    Name:     cfg.Database.Name,
    SSLMode:  cfg.Database.SSLMode,
	})
	defer db.Close()

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
