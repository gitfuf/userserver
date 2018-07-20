package usecases

import (
	"log"
	"net/http"

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

	//use standart mux for simple requests (actually better to use httprouter, ut for test purpose will stay so )
	stdRouter := http.NewServeMux()
	stdRouter.HandleFunc("/shutdown", app.serverShutdown)

	serverMux := http.NewServeMux()
	serverMux.Handle("/user/", gorillaRouter)
	serverMux.Handle("/user", httpRouter)
	serverMux.Handle("/shutdown", stdRouter)
	app.Server.Handler = serverMux
}
