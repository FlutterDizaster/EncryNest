package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/FlutterDizaster/EncryNest/internal/server"
	"github.com/FlutterDizaster/EncryNest/pkg/configloader"
)

func main() {
	// Gracefull shutdown with SIGINT and SIGTERM
	ctx, cancle := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancle()

	var settings server.Settings

	err := configloader.LoadConfig(&settings)
	if err != nil {
		slog.Error("Error while loading config", slog.Any("err", err))
		return
	}

	srv := server.NewServer(settings)

	err = srv.Init(ctx)
	if err != nil {
		slog.Error("Error while initializing server", slog.Any("err", err))
		return
	}

	err = srv.Run(ctx)
	if err != nil {
		slog.Error("Error while running server", slog.Any("err", err))
	}
}
