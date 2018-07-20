package usecases

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
)

func (app *ServerApp) initRESTApi() {
	log.Println("ServerApp:initRESTApi")

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

	serverMux := http.NewServeMux()
	serverMux.Handle("/user/", gorillaRouter)
	serverMux.Handle("/user", httpRouter)
	serverMux.Handle("/shutdown", stdRouter)

	//set middleware
	serverHandler := requestsLogMiddleware(serverMux)
	serverHandler = panicMiddleware(serverHandler)

	app.Server.Handler = serverHandler
}

func panicMiddleware(mw http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//log.Println("panicMiddleware ", r.URL.Path)
		defer func() {
			if err := recover(); err != nil {
				log.Printf("recovered from error: %v ", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		mw.ServeHTTP(w, r)
	})
}

func requestsLogMiddleware(mw http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("requestsLogMiddleware ", r.URL.Path)
		start := time.Now()
		mw.ServeHTTP(w, r)
		log.Printf("[%s] %s, %s %s\n", r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))
	})
}
