package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/SemenShakhray/sso/internal/app"
	"github.com/SemenShakhray/sso/internal/config"
	"github.com/SemenShakhray/sso/internal/lib/logger"
)

func main() {
	// TODO: init config
	cfg := config.MustLoad()

	//TODO: init logger
	log := logger.SetupLogger(cfg.Env)

	log.Info("starting application",
		slog.Any("cfg", cfg),
	)

	//TODO: init app
	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	//TODO: run grpc-server
	go func() {
		application.GRPCSrv.MustRun()
	}()

	//TODO: Graceful shutdown
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	stop := <-sigint

	application.GRPCSrv.Stop()

	log.Info("stopped application",
		slog.String("received signal", stop.String()))
}
