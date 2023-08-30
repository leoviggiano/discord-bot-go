package main

import (
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"

	"bot_test/bot"
	"bot_test/lib/log"
	"bot_test/service"
)

func main() {
	logOptions := make([]log.Option, 0)
	logOptions = append(logOptions, log.WithTextFormatter())
	log := log.NewLogger(logrus.InfoLevel, logOptions...)

	services := service.GetAll(log)
	if err := bot.StartBot(services, log); err != nil {
		log.Fatal(err)
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
