package usecases

import (
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

func (app *ServerApp) initRESTApi() {
	log.Debug("ServerApp:initRESTApi")

	//Create gorilla mux for complex requests
	gorillaRouter := mux.NewRouter()
	gorillaRouter.HandleFunc("/user/{id:[0-9]+}", app.getUser).Methods("GET")
	gorillaRouter.HandleFunc("/user/{id:[0-9]+}", app.updateUser).Methods("PUT")
	gorillaRouter.HandleFunc("/user/{id:[0-9]+}", app.deleteUser).Methods("DELETE")

	//Create httproute for /user which not need id check and will be faster
	httpRouter := httprouter.New()
	httpRouter.POST("/user", app.newUser)
	//gorillaRouter.HandleFunc("/user", app.newUser).Methods("POST")

	//use standart mux for simple requests (actually better to use httprouter, but for test purpose will stay so )
	stdRouter := http.NewServeMux()
	stdRouter.HandleFunc("/shutdown", app.serverShutdown)
	stdRouter.HandleFunc("/panic", app.panicHappen)

	serverMux := http.NewServeMux()
	serverMux.Handle("/user/", gorillaRouter)
	serverMux.Handle("/user", httpRouter)
	serverMux.Handle("/shutdown", stdRouter)
	serverMux.Handle("/panic", stdRouter)

	//set middleware
	serverHandler := requestsLogMiddleware(serverMux)
	serverHandler = panicMiddleware(serverHandler)

	app.Server.Handler = serverHandler
}

func panicMiddleware(mw http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("recovered from error: %v \n", err)
				stack := debug.Stack()
				log.Errorln(string(stack))
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		mw.ServeHTTP(w, r)
	})
}

func requestsLogMiddleware(mw http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknown"
		}

		mw.ServeHTTP(w, r)
		log.WithFields(log.Fields{
			"url_path":    r.URL.Path,
			"method":      r.Method,
			"hostname":    hostname,
			"remote_addr": r.RemoteAddr,
			"work_time":   time.Since(start),
		}).Info("New request:")
	})
}
