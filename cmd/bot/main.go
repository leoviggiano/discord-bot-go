package main

import (
	"log/slog"
	"os"
	"os/signal"

	"megumin/config"
	"megumin/pkg/bot"
	"megumin/pkg/lib/sql"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if config.Environment() != config.Production {
		logger.Info("Starting bot in development mode")
	}

	db, err := sql.NewSQLXConnection(sql.Settings{
		Conn: config.DatabaseConnString(),
	}, logger)

	if err != nil {
		logger.Error("Error database connection", "error", err)
	}

	if err := bot.StartBot(config.DiscordToken(), db, logger); err != nil {
		logger.Error("Error starting bot", "error", err)
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		close(idleConnsClosed)
	}()

	<-idleConnsClosed
}
