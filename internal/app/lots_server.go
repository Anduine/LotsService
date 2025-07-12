package app

import (
	"log/slog"
	"server/internal/delivery/http_handlers"
	"server/internal/repository"
	"server/internal/server"
	"server/internal/service"
	"server/pkg/database"
	"time"
)

func Run(log *slog.Logger, port string, timeout time.Duration, dbConn string) {
	db := database.NewPostgresConnection(dbConn)

	repo := repository.NewPostgresLotsRepo(db)
	lotsService := service.NewLotsService(repo)
	lotsHandler := http_handlers.NewLotsHandler(lotsService)

	handler := server.NewRouter(lotsHandler)

	server.StartServer(log, handler, port, timeout)
}