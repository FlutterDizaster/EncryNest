package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/FlutterDizaster/EncryNest/internal/server"
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

	settings := server.Settings{
		Addr:      "",
		Port:      "50555",
		JWTSecret: "secret",
	}
	srv := server.NewServer(settings)

	err := srv.Run(ctx)
	if err != nil {
		slog.Error("Error while running server", slog.Any("err", err))
	}
}
