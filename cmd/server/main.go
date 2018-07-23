//Copyright © 2018 Fuf
//Test server app for practice http REST API and DB logic patterns
package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/gitfuf/userserver/config"
	"github.com/gitfuf/userserver/repository"
	"github.com/gitfuf/userserver/repository/handlers"
	"github.com/gitfuf/userserver/usecases"
)

func main() {
	logfile := config.InitLog()

	//close log file at the app exit
	defer func() {
		fmt.Println("Close log file: ", logfile.Name())
		logfile.Sync()
		logfile.Close()
	}()

	err := config.InitConfiguration("pro", "")
	if err != nil {
		log.Fatal("InitConfiguration err=", err)
	}

	dbRepo, err := setupDBRepo()
	if err != nil {
		log.Error("Setup new DB repository err: ", err)
		return
	}
	defer dbRepo.CloseDB()

	server, err := usecases.NewServer(dbRepo, port())
	if err != nil {
		log.Error("Create new server err: ", err)
	}

	done := make(chan bool)
	go func() {
		log.Debug("Run ListenAndServe()")
		err := server.ListenAndServe()
		if err != nil {
			log.Warnf("Listen and serve: %v", err)
		}
		done <- true
	}()

	//wait shutdown
	server.WaitShutdown()

	<-done
	log.Info("Server was graceful shutdown")

}

func port() string {
	ret := os.Getenv("HTTP_PORT")
	if ret != "" {
		return ret
	}
	ret = ":" + config.HttpPort()
	log.Debug("port = ", ret)
	return ret
}

func setupDBRepo() (usecases.DBRepository, error) {
	switch config.DBDriver() {
	case "postgres":
		postgresHandler, err := handlers.NewPostgresHandler(config.DBConnString())
		if err != nil {
			log.Warn("setupDB postgres error=", err)
			return nil, err
		}
		postgresRepo := repository.NewPostgresRepository(postgresHandler)
		return postgresRepo, nil
	case "mssql":
		msHandler, err := handlers.NewMssqlHandler(config.DBConnString())
		if err != nil {
			log.Warn("setupDB mssql error=", err)
			return nil, err
		}
		msRepo := repository.NewMssqlRepository(msHandler)
		return msRepo, nil
	case "mysql":
		myHandler, err := handlers.NewMysqlHandler(config.DBConnString())
		if err != nil {
			log.Warn("setupDB mysql error=", err)
			return nil, err
		}
		myRepo := repository.NewMysqlRepository(myHandler)
		return myRepo, nil
	default:
		log.Warn("setupDB: unknow driver")
	}
	return nil, nil

}
