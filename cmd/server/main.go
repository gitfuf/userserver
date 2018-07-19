//Copyright © 2018 Fuf
//Test server app for practice http REST API and DB logic patterns
package main

import (
	//	"fmt"
	"log"
	"os"

	"github.com/gitfuf/userserver/config"
	"github.com/gitfuf/userserver/repository"
	"github.com/gitfuf/userserver/repository/handlers"
	"github.com/gitfuf/userserver/usecases"
)

func main() {
	err := config.InitVars("pro", "")
	if err != nil {
		log.Println("InitVars err=", err)
		os.Exit(1)
	}

	dbRepo, err := setupDBRepo()
	if err != nil {
		log.Fatal(err)
	}
	defer dbRepo.CloseDB()

	server, err := usecases.NewServer(dbRepo, port())
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)
	go func() {
		log.Println("Run ListenAndServe()")
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("Listen and serve: %v", err)
		}
		done <- true
	}()

	//wait shutdown
	server.WaitShutdown()

	<-done
	log.Println("Server was graceful shutdown")

}

func port() string {
	ret := os.Getenv("HTTP_PORT")
	if ret != "" {
		return ret
	}
	ret = ":" + config.HttpPort()
	log.Println("port = ", ret)
	return ret
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
