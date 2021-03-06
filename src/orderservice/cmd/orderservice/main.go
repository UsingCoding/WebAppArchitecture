package main

import (
	"context"
	"net/http"
	"orderservice/pkg/orderservice/domain/service"
	"orderservice/pkg/orderservice/infrastructure/mysql/query"
	"orderservice/pkg/orderservice/infrastructure/mysql/repository"
	"orderservice/pkg/orderservice/infrastructure/transport"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

var appID = "UNKNOWN"

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	config, err := parseEnv()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sqlx.Open(config.DatabaseDriver, config.DSN)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	killSignalChan := getKillSignalChan()

	menuItemRepo := repository.NewMenuItemRepository(db)
	menuItemService := service.NewMenuItemService(menuItemRepo)
	menuItemQueryService := query.NewMenuItemQueryService(db)

	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo, menuItemRepo)
	orderQueryService := query.NewOrderQueryService(db, menuItemQueryService)

	server := startServer(config.ServeHTTPAddress, transport.NewServer(
		orderService,
		orderQueryService,
		menuItemService,
		menuItemQueryService,
	))

	waitForKillSignal(killSignalChan)
	err = server.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("failed to stop server: %s", err)
	}
}

func startServer(addr string, srv *transport.Server) *http.Server {
	router := transport.NewRouter(srv)
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
