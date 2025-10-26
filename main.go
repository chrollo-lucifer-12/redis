package main

import (
	"log/slog"
	"net"
	"os"

	"github.com/chrollo-lucifer-12/redis/server"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	address := ":6379"
	logger.Info("starting server", slog.String("address", address))
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error(
			"cannot start tcp server",
			slog.String("address", address),
			slog.String("err", err.Error()),
		)
		os.Exit(-1)
	}
	s := server.NewServer(listener, logger)
	go func() {
		if err := s.Start(); err != nil {
			logger.Error("server error", slog.String("err", err.Error()))
			os.Exit(1)
		}
	}()
	select {}
}
