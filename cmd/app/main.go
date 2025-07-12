package main

import (
	"log/slog"
	"os"

	"server/internal/app"
	"server/internal/config"

	slogplus "server/internal/lib/logger"
)

func main() {
	config := config.MustLoadConfig()

	log := setupPlusSlog()

	app.Run(log, config.Port, config.Timeout, config.DBConnector)
}


func setupPlusSlog() *slog.Logger {
	opts := slogplus.PlusHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPlusHandler(os.Stdout)

	return slog.New(handler)
}