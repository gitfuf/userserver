package usecases

//setup REST API
func (app *ServerApp) initRoutes() {
	app.Router.HandleFunc("/user", app.newUser).Methods("POST")
	//	app.Router.HandleFunc("/user/emails", app.getEmails).Methods("GET")
	app.Router.HandleFunc("/user/{id:[0-9]+}", app.getUser).Methods("GET")
	app.Router.HandleFunc("/user/{id:[0-9]+}", app.updateUser).Methods("PUT")
	app.Router.HandleFunc("/user/{id:[0-9]+}", app.deleteUser).Methods("DELETE")

}
