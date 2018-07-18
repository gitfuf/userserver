//Copyright © 2018 Fuf
package usecases

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type ServerApp struct {
	Router *mux.Router
	DBRepo DBRepository

	//for gracefulShutdown
	http.Server
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
	srv.initRouter()

	return srv, nil
}

func (app *ServerApp) initRouter() {
	log.Println("ServerApp:initRouter()")
	router := mux.NewRouter()
	app.Router = router
	app.Server.Handler = router
	app.initRoutes()
}

func (app *ServerApp) WaitShutdown() {
	irqSig := make(chan os.Signal, 1)
	signal.Notify(irqSig, syscall.SIGINT, syscall.SIGTERM)

	//Wait interrupt or shutdown request through /shutdown
	select {
	case sig := <-irqSig:
		log.Printf("Shutdown request (signal: %v)", sig)
	case sig := <-app.shutdownC:
		log.Printf("Shutdown request (/shutdown %v)", sig)
	}

	log.Printf("Stoping http server ...")

	//Create shutdown context with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//shutdown the server
	err := app.Shutdown(ctx)
	if err != nil {
		log.Printf("Shutdown request error: %v", err)
	}
}
