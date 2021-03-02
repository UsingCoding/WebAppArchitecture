package main

import (
	"context"
	"net/http"
	"orderservice/pkg/orderservice/infrastructure/transport"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile("dev.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
		defer file.Close()
	}

	killSignalChan := getKillSignalChan()

	addr := ":8000"
	server := startServer(addr)
	waitForKillSignal(killSignalChan)
	err = server.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("failed to stop server: %s", err)
	}
}

func startServer(addr string) *http.Server {
	router := transport.NewRouter()
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		log.WithFields(log.Fields{"addr": addr}).Info("starting the server")
		log.Info(server.ListenAndServe())
	}()

	return server
}

func getKillSignalChan() chan os.Signal {
	osKillSignalChan := make(chan os.Signal, 1)
	signal.Notify(osKillSignalChan, os.Interrupt, syscall.SIGTERM)
	return osKillSignalChan
}

func waitForKillSignal(killSignalChan <-chan os.Signal) {
	killSignal := <-killSignalChan
	switch killSignal {
	case os.Interrupt:
		log.Info("got SIGINT")
	case syscall.SIGTERM:
		log.Info("got SIGTERM")
	}
}
