//Copyright © 2018 Fuf
//Test server app for practice http REST API and DB logic patterns
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gitfuf/userserver/config"
	"github.com/gitfuf/userserver/repository"
	"github.com/gitfuf/userserver/repository/handlers"
	"github.com/gitfuf/userserver/usecases"
)

func main() {
	//log.Println("PG_HOST=", os.Getenv("PG_HOST"))
	err := config.InitVars("pro", "")
	if err != nil {
		fmt.Println("err=", err)
		os.Exit(3)
	}

	app := usecases.ServerApp{}

	dbRepo, err := setupDBRepo()
	if err != nil {
		log.Fatal(err)
	}
	app.DBRepo = dbRepo
	app.InitRouter()
	app.Run(":8080")

}

func init() {

}

func setupDBRepo() (usecases.DBRepository, error) {
	switch config.DBDriver() {
	case "postgres":
		postgresHandler, err := handlers.NewPostgresHandler(config.DBConnString())
		if err != nil {
			log.Println("setupDB postgres error=", err)
			return nil, err
		}
		postgresRepo := repository.NewPostgresRepository(postgresHandler)
		return postgresRepo, nil
	case "mssql":
		msHandler, err := handlers.NewMssqlHandler(config.DBConnString())
		if err != nil {
			log.Println("setupDB mssql error=", err)
			return nil, err
		}
		msRepo := repository.NewMssqlRepository(msHandler)
		return msRepo, nil
	case "mysql":
		myHandler, err := handlers.NewMysqlHandler(config.DBConnString())
		if err != nil {
			log.Println("setupDB mysql error=", err)
			return nil, err
		}
		myRepo := repository.NewMysqlRepository(myHandler)
		return myRepo, nil
	default:
		log.Println("setupDB: unknow driver")
	}
	return nil, nil

}
