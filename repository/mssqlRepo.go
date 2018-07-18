package repository

import (
	"errors"
	"log"

	"github.com/gitfuf/userserver/repository/handlers"
	"github.com/gitfuf/userserver/usecases"
)

type MssqlRepo struct {
	usecases.DBRepository
	msHandler *handlers.MsHandler
}

func NewMssqlRepository(msHandler *handlers.MsHandler) usecases.DBRepository {
	log.Println("call NewMssqlRepository")
	msRepo := new(MssqlRepo)
	msRepo.msHandler = msHandler
	return msRepo
}

func (msRepo *MssqlRepo) AddUser(u *usecases.User) error {
	log.Println("MssqlRepo: AddUser:begin", u)

	uM := createModelUser(*u)
	err := msRepo.msHandler.InsertUser(&uM)
	if err != nil {
		log.Println("MssqlRepo:AddUser err=", err)
		return err
	}
	//TODO maybe convert func
	u.ID = uM.ID
	log.Println("MssqlRepo:AddUser:success =", u)
	return nil
}

func (msRepo *MssqlRepo) UpdateUser(u usecases.User) error {
	log.Println("MssqlRepo:UpdateUser:begin user=", u)
	uM := createModelUser(u)
	err := msRepo.msHandler.UpdateUser(uM)
	if err != nil {
		log.Println("MssqlRepo:UpdateUser:error=", err)
		return err
	}
	log.Println("MssqlRepo:UpdateUser:successful")
	return nil
}
func (msRepo *MssqlRepo) DeleteUser(id int64) error {
	log.Println("MssqlRepo:DeleteUser:begin id=", id)

	err := msRepo.msHandler.DeleteUser(id)
	if err != nil {
		log.Println("MssqlRepo:DeleteUser:error=", err)
		return err
	}
	log.Println("MssqlRepo:DeleteUser:successful")
	return nil
}
func (msRepo *MssqlRepo) GetUserInfo(id int64) (usecases.User, error) {
	log.Println("MssqlRepo:GetUserInfo begin id=", id)
	uM, err := msRepo.msHandler.GetUser(id)
	u := createUcUser(uM)
	log.Printf("MssqlRepo:GetUserInfo result user=%v, err=%v\n", u, err)
	return u, err

}

func (msRepo *MssqlRepo) CloseDB() {
	if err := msRepo.msHandler.CloseDB(); err == nil {
		log.Printf("MssqlRepo: CloseDB() successful")
	} else {
		log.Printf("MssqlRepo: CloseDB() err %v\n", err)
	}
}

func (msRepo *MssqlRepo) CreateTable(table string) error {
	log.Println("MssqlRepo:CreateTable ", table)
	switch table {
	case "users":
		return msRepo.msHandler.CreateUserTable()
	default:
		break
	}
	return errors.New("unknown table:" + table)
}

func (msRepo *MssqlRepo) ClearTable(table string) error {
	log.Println("MssqlRepo:ClearTable ", table)
	switch table {
	case "users":
		return msRepo.msHandler.ClearUserTable()
	default:
		break
	}
	return errors.New("unknown table:" + table)
}
