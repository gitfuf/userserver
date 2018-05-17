//Copyright © 2018 Fuf
package usecases

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ServerApp struct {
	Router *mux.Router
	DBRepo DBRepository
}

func (app *ServerApp) InitRouter() {
	app.Router = mux.NewRouter()
	app.initRoutes()
}

//start http server with router
func (app *ServerApp) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, app.Router))
}

func (app *ServerApp) newUser(w http.ResponseWriter, r *http.Request) {
	log.Println("ServerApp:newUser route")
	u := User{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request data")
		return
	}
	defer r.Body.Close()

	//fmt.Println(err)
	log.Println("ServerApp:call DBRepo.AddUser:", u)
	err := app.DBRepo.AddUser(&u)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, u)
}

func (app *ServerApp) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	u, err := app.DBRepo.GetUserInfo(int64(id))
	if err != nil {
		switch err.Error() {
		case "haven't found":
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	log.Println("ServerApp:getUser success=", u)
	respondWithJSON(w, http.StatusOK, u)
}

func (app *ServerApp) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var u User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest data")
		return
	}
	defer r.Body.Close()
	u.ID = int64(id)

	if err := app.DBRepo.UpdateUser(u); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (app *ServerApp) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	if err := app.DBRepo.DeleteUser(int64(id)); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

/*
func (app *ServerApp) getEmails(w http.ResponseWriter, r *http.Request) {

	emails, err := app.DB.SelectAllEmails()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, emails)
}
*/
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
