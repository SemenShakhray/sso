package main

import (
	"log/slog"

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
	application.GRPCSrv.MustRun()
}
