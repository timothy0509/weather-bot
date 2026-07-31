package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"weather-bot/internal/bot"
	"weather-bot/internal/config"
	"weather-bot/internal/db"
	"weather-bot/internal/monitor"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", slog.Any("err", err))
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}))

	logger.Info("opening database", slog.String("path", cfg.DatabasePath))
	database, err := db.Open(cfg.DatabasePath)
	if err != nil {
		logger.Error("failed to open database", slog.Any("err", err))
		os.Exit(1)
	}
	defer database.Close()

	b, err := bot.New(cfg, database, logger)
	if err != nil {
		logger.Error("failed to create bot", slog.Any("err", err))
		os.Exit(1)
	}

	b.RegisterHandlers()

	if err := b.Session.Open(); err != nil {
		logger.Error("failed to open discord session", slog.Any("err", err))
		os.Exit(1)
	}

	if err := b.RegisterCommands(); err != nil {
		logger.Error("failed to register commands", slog.Any("err", err))
		os.Exit(1)
	}

	mon := monitor.New(b.Session, database, b.HKO, logger, monitor.Intervals{
		Warning: cfg.WarningInterval,
		Tips:    cfg.TipsInterval,
		Status:  cfg.StatusInterval,
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mon.Start(ctx)

	logger.Info("bot started", slog.String("guild_id", cfg.GuildID))

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	logger.Info("shutting down")
	cancel()
	mon.Wait()
	_ = b.Close()
}
