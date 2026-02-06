package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"remora/internal/agent"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	slog.SetDefault(logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rebalanceStop, err := agent.StartCron(ctx, logger, true)
	if err != nil {
		logger.Error("rebalance cron start failed", slog.Any("error", err))
		os.Exit(1)
	}

	if rebalanceStop != nil {
		defer rebalanceStop()
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM)

	<-interrupt
	logger.Info("shutting down...")
	cancel()
}
