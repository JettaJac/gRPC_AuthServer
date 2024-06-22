package app

import (
	"log/slog"
	"time"

	"os"
	grpcapp "sso/internal/app/grpc"
	"sso/internal/config"
	"sso/internal/lib/logger"
	"sso/internal/services/auth"
	"sso/internal/storage/posgre"
	// "sso/internal/storage/sqlite"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New( // возможно чтоб просто принимал конфиг
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
	config *config.Config,
) *App {
	// storage, err := sqlite.New(storagePath)
	// if err != nil {
	// 	panic(err)
	// }

	storage, err := posgre.New(config.DatabaseURL)
	if err != nil {
		log.Error("app.Run: Failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	log.Info("Сonnected to the database", slog.String("env", config.Env))
	defer storage.CloseDB()

	authService := auth.New(log, storage, storage, storage, tokenTTL) //переделать на 1 стораж
	log.Info("Starting server", slog.String("address", string(config.GRPC.Port)))
	grpcApp := grpcapp.New(log, authService, grpcPort)
	log.Info("Starting app", slog.String("address", string(config.GRPC.Port))) //!!  adress will delete
	return &App{
		GRPCSrv: grpcApp,
	}
}
