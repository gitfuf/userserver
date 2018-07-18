package repository

import (
	"errors"
	"log"

	"github.com/gitfuf/userserver/repository/handlers"
	"github.com/gitfuf/userserver/usecases"
)

type PostgresRepo struct {
	usecases.DBRepository
	pgHandler *handlers.PostgresHandler
}

func NewPostgresRepository(pgHandler *handlers.PostgresHandler) usecases.DBRepository {
	log.Println("call NewPostgresRepository")
	pgRepo := new(PostgresRepo)
	pgRepo.pgHandler = pgHandler
	return pgRepo
}

func (pgRepo *PostgresRepo) GetUserInfo(id int64) (usecases.User, error) {
	log.Println("PostgresRepo:GetUserInfo begin id=", id)
	uM, err := pgRepo.pgHandler.GetUser(id)
	u := createUcUser(uM)
	log.Printf("PostgresRepo:GetUserInfo result user=%v, err=%v\n", u, err)
	return u, err

}

func (pgRepo *PostgresRepo) AddUser(u *usecases.User) error {
	log.Println("PostgresRepo: AddUser:begin", u)

	uM := createModelUser(*u)
	err := pgRepo.pgHandler.InsertUser(&uM)
	if err != nil {
		log.Println("PostgresRepo:AddUser err=", err)
		return err
	}
	//TODO maybe convert func
	u.ID = uM.ID
	log.Println("PostgresRepo:AddUser:success =", u)
	return nil
}

func (pgRepo *PostgresRepo) UpdateUser(u usecases.User) error {
	log.Println("PostgresRepo:UpdateUser:begin user=", u)
	uM := createModelUser(u)
	err := pgRepo.pgHandler.UpdateUser(uM)
	if err != nil {
		log.Println("PostgresRepo:UpdateUser:error=", err)
		return err
	}
	log.Println("PostgresRepo:UpdateUser:successful")
	return nil
}
func (pgRepo *PostgresRepo) DeleteUser(id int64) error {
	log.Println("PostgresRepo:DeleteUser:begin id=", id)

	err := pgRepo.pgHandler.DeleteUser(id)
	if err != nil {
		log.Println("PostgresRepo:DeleteUser:error=", err)
		return err
	}
	log.Println("PostgresRepo:DeleteUser:successful")
	return nil
}

func (pgRepo *PostgresRepo) CloseDB() {
	if err := pgRepo.pgHandler.CloseDB(); err == nil {
		log.Printf("PostgresRepo: CloseDB() successful")
	} else {
		log.Printf("PostgresRepo: CloseDB() err %v\n", err)
	}

}

func (pgRepo *PostgresRepo) CreateTable(table string) error {
	log.Println("PostgresRepo:CreateTable ", table)
	switch table {
	case "users":
		return pgRepo.pgHandler.CreateUserTable()
	default:
		break
	}
	return errors.New("unknown table:" + table)
}

func (pgRepo *PostgresRepo) ClearTable(table string) error {
	switch table {
	case "users":
		return pgRepo.pgHandler.ClearUserTable()
	default:
		break
	}
	return errors.New("unknown table:" + table)
}
