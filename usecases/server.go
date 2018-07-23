//Copyright © 2018 Fuf
package usecases

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

type ServerApp struct {
	DBRepo DBRepository
	http.Server

	//for graceful shutdown
	shutdownC chan bool
	reqCount  uint32
}

func NewServer(db DBRepository, port string) (*ServerApp, error) {
	srv := &ServerApp{
		Server: http.Server{
			Addr: ":8080",
		},
		shutdownC: make(chan bool),
	}
	srv.DBRepo = db
	srv.initRESTApi()

	return srv, nil
}

func (app *ServerApp) WaitShutdown() {
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)

	//Wait interrupt or shutdown request through /shutdown
	select {
	case sig := <-sigC:
		log.Debugf("Shutdown request (signal: %v)", sig)
	case sig := <-app.shutdownC:
		log.Debugf("Shutdown request (/shutdown %v)", sig)
	}

	log.Info("Stoping http server ...")

	//Create shutdown context with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//shutdown the server
	err := app.Shutdown(ctx)
	if err != nil {
		log.Warn("Shutdown request error: %v", err)
	}
}
