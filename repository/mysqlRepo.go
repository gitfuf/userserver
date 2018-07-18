package repository

import (
	"errors"
	"log"

	"github.com/gitfuf/userserver/repository/handlers"
	"github.com/gitfuf/userserver/usecases"
)

type MysqlRepo struct {
	usecases.DBRepository
	myHandler *handlers.MysqlHandler
}

func NewMysqlRepository(myHandler *handlers.MysqlHandler) usecases.DBRepository {
	log.Println("call NewMysqlRepository")
	myRepo := new(MysqlRepo)
	myRepo.myHandler = myHandler
	return myRepo
}

func (myRepo *MysqlRepo) GetUserInfo(id int64) (usecases.User, error) {
	log.Println("MysqlRepo:GetUserInfo begin id=", id)
	uM, err := myRepo.myHandler.GetUser(id)
	u := createUcUser(uM)
	log.Printf("MysqlRepo:GetUserInfo result user=%v, err=%v\n", u, err)
	return u, err

}

func (myRepo *MysqlRepo) AddUser(u *usecases.User) error {
	log.Println("MysqlRepo: AddUser:begin", u)

	uM := createModelUser(*u)
	err := myRepo.myHandler.InsertUser(&uM)
	if err != nil {
		log.Println("MysqlRepo:AddUser err=", err)
		return err
	}
	//TODO maybe convert func
	u.ID = uM.ID
	log.Println("MysqlRepo:AddUser:success =", u)
	return nil
}

func (myRepo *MysqlRepo) UpdateUser(u usecases.User) error {
	log.Println("MysqlRepo:UpdateUser:begin user=", u)
	uM := createModelUser(u)
	err := myRepo.myHandler.UpdateUser(uM)
	if err != nil {
		log.Println("MysqlRepo:UpdateUser:error=", err)
		return err
	}
	log.Println("MysqlRepo:UpdateUser:successful")
	return nil
}
func (myRepo *MysqlRepo) DeleteUser(id int64) error {
	log.Println("MysqlRepo:DeleteUser:begin id=", id)

	err := myRepo.myHandler.DeleteUser(id)
	if err != nil {
		log.Println("MysqlRepo:DeleteUser:error=", err)
		return err
	}
	log.Println("MysqlRepo:DeleteUser:successful")
	return nil
}

func (myRepo *MysqlRepo) CreateTable(table string) error {
	log.Println("MysqlRepo:CreateTable ", table)
	switch table {
	case "users":
		return myRepo.myHandler.CreateUserTable()
	default:
		break
	}
	return errors.New("unknown table:" + table)
}

func (myRepo *MysqlRepo) ClearTable(table string) error {
	switch table {
	case "users":
		return myRepo.myHandler.ClearUserTable()
	default:
		break
	}
	return errors.New("unknown table:" + table)
}
