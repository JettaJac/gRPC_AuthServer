package main

// protoc -I proto --go_out=plugins=grpc:cmd proto/sso/sso.proto
// protoc -I proto proto/sso/sso.proto --go_out=./gen/go --go_opt=paths=source_relative --go-grpc_out=./gen/go/ --go-grpc_opt=paths=source_relative
// go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
// go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"sso/internal/app"
	"sso/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	fmt.Println("Start app")

	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("starting application", slog.String("env", cfg.Env))
	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL, cfg)
	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop
	log.Info("stopping signal", slog.String("signal", sign.String()))
	application.GRPCSrv.Stop()
	log.Info("application stopped")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
