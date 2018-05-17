package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gitfuf/userserver/config"
	//"github.com/gitfuf/userserver/repository"
	//"github.com/gitfuf/userserver/repository/handlers"
	"github.com/gitfuf/userserver/usecases"
)

const userTable = "users"

var app usecases.ServerApp

func TestMain(m *testing.M) {

	tests := []struct {
		name   string
		driver string
	}{
		{
			"postgres tests",
			"postgres",
		},
		{
			"mssql tests",
			"mssql",
		},
	}
	for _, tt := range tests {
		fmt.Println("Start tests: ", tt.driver)
		err := config.InitVars("test", tt.driver)
		if err != nil {
			os.Exit(1)
		}
		app = usecases.ServerApp{}

		dbRepo, err := setupDBRepo()
		if err != nil {
			log.Fatal(err)
		}
		app.DBRepo = dbRepo
		app.InitRouter()

		ensureTableExists(userTable)

		m.Run()

		clearTable(userTable)
	}

	//os.Exit(code)
}

func TestUser_New(t *testing.T) {

	tests := []struct {
		name               string
		data               []byte
		wantErr            bool
		httpExpectedStatus int
	}{
		{
			"add new user successfull",
			[]byte(`{"age":24,"first_name":"Ma","last_name":"So","email":"fuf@fu.com"}`),
			false,
			http.StatusCreated,
		},
		{
			"add new user fail no data",
			[]byte(``),
			true,
			http.StatusBadRequest,
		},
	}

	clear := func() {
		clearTable(userTable)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clear()
			defer clear()

			req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(tt.data))
			response := executeRequest(req)

			if !tt.wantErr {
				var m map[string]interface{}
				json.Unmarshal(response.Body.Bytes(), &m)

				if m["age"] != 24.0 {
					t.Errorf("Expected user age to be '24'. Got '%v'", m["age"])
				}

				if m["first_name"] != "Ma" {
					t.Errorf("Expected first name to be 'Ma'. Got '%v'", m["first_name"])
				}

				if m["last_name"] != "So" {
					t.Errorf("Expected last name to be 'So'. Got '%v'", m["last_name"])
				}

				if m["email"] != "fuf@fu.com" {
					t.Errorf("Expected last name to be 'fuf@fu.com'. Got '%v'", m["email"])
				}

				if m["id"] != 1.0 {
					t.Errorf("Expected user ID to be '1'. Got '%v'", m["id"])
				}
			}
			if tt.httpExpectedStatus != response.Code {
				t.Errorf("Expected response code %d. Got %d\n", tt.httpExpectedStatus, response.Code)
			}
		})
	}
}

func TestUser_Update(t *testing.T) {

	tests := []struct {
		name               string
		data               []byte
		wantErr            bool
		httpExpectedStatus int
	}{
		{
			"update user successfull",
			[]byte(`{"age":24,"first_name":"Ma","last_name":"So","email":"fuf@fu.com"}`),
			false,
			http.StatusOK,
		},
		{
			"update user fail no data",
			[]byte(``),
			true,
			http.StatusBadRequest,
		},
	}

	clear := func() {
		clearTable(userTable)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clear()
			defer clear()

			u := usecases.User{
				Age:       11,
				FirstName: "Ka",
				LastName:  "Po",
				Email:     "h@h.com",
			}
			app.DBRepo.AddUser(&u)
			req, _ := http.NewRequest("GET", "/user/1", nil)
			response := executeRequest(req)
			var originalUser map[string]interface{}
			json.Unmarshal(response.Body.Bytes(), &originalUser)

			req, _ = http.NewRequest("PUT", "/user/1", bytes.NewBuffer(tt.data))
			response = executeRequest(req)

			if !tt.wantErr {
				var m map[string]interface{}
				json.Unmarshal(response.Body.Bytes(), &m)

				if m["id"] != originalUser["id"] {
					t.Errorf("Expected the id to remain the same (%v). Got %v", originalUser["id"], m["id"])
				}

				if m["first_name"] == originalUser["first_name"] {
					t.Errorf("Expected the first_name to change from '%v' to '%v'. Got '%v'", originalUser["first_name"], m["first_name"], m["first_name"])
				}

				if m["last_name"] == originalUser["last_name"] {
					t.Errorf("Expected the last_name to change from '%v' to '%v'. Got '%v'", originalUser["last_name"], m["last_name"], m["last_name"])
				}

				if m["email"] == originalUser["email"] {
					t.Errorf("Expected the email to change from '%v' to '%v'. Got '%v'", originalUser["email"], m["email"], m["email"])
				}

				if m["age"] == originalUser["age"] {
					t.Errorf("Expected the age to change from '%v' to '%v'. Got '%v'", originalUser["age"], m["age"], m["age"])
				}
			}
			if tt.httpExpectedStatus != response.Code {
				t.Errorf("Expected response code %d. Got %d\n", tt.httpExpectedStatus, response.Code)
			}
		})
	}
}

func TestUser_Delete(t *testing.T) {

	tests := []struct {
		name string
		url  string
		//data               []byte
		wantErr            bool
		httpExpectedStatus int
	}{
		{
			"delete user successfull",
			"/user/1",
			false,
			http.StatusOK,
		},
		{
			"delete user fail bad id",
			"/user/",
			true,
			http.StatusNotFound,
		},
		{
			"delete user fail non existing id",
			"/user/23",
			true,
			http.StatusInternalServerError,
		},
	}

	clear := func() {
		clearTable(userTable)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clear()
			defer clear()

			u := usecases.User{
				Age:       11,
				FirstName: "Ka",
				LastName:  "Po",
				Email:     "h@h.com",
			}
			app.DBRepo.AddUser(&u)
			data := []byte("")
			req, _ := http.NewRequest("DELETE", tt.url, bytes.NewBuffer(data))
			response := executeRequest(req)

			if tt.httpExpectedStatus != response.Code {
				t.Errorf("Expected response code %d. Got %d\n", tt.httpExpectedStatus, response.Code)
			}
		})
	}
}

func TestUser_Get(t *testing.T) {

	tests := []struct {
		name               string
		url                string
		wantErr            bool
		httpExpectedStatus int
	}{
		{
			"get user successfull",
			"/user/1",
			false,
			http.StatusOK,
		},
		{
			"get user fail bad id",
			"/user/",
			true,
			http.StatusNotFound,
		},
		{
			"get user fail non existing id",
			"/user/23",
			true,
			http.StatusNotFound,
		},
	}

	clear := func() {
		clearTable(userTable)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clear()
			defer clear()

			u := usecases.User{
				Age:       11,
				FirstName: "Ka",
				LastName:  "Po",
				Email:     "h@h.com",
			}
			app.DBRepo.AddUser(&u)
			data := []byte("")
			req, _ := http.NewRequest("GET", tt.url, bytes.NewBuffer(data))
			response := executeRequest(req)

			if tt.httpExpectedStatus != response.Code {
				t.Errorf("Expected response code %d. Got %d\n", tt.httpExpectedStatus, response.Code)
			}
		})
	}
}

// http
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	return rr
}

/*
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
*/
//DB
func ensureTableExists(table string) {
	err := app.DBRepo.CreateTable(table)
	if err != nil {
		log.Fatal(err)
	}
}

func clearTable(table string) {

	if err := app.DBRepo.ClearTable(table); err != nil {
		log.Fatal(err)
	}
}
