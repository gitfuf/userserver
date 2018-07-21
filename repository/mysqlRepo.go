package repository

import (
	"errors"

	"github.com/gitfuf/userserver/repository/handlers"
	"github.com/gitfuf/userserver/usecases"
	log "github.com/sirupsen/logrus"
)

type MysqlRepo struct {
	usecases.DBRepository
	myHandler *handlers.MysqlHandler
}

func NewMysqlRepository(myHandler *handlers.MysqlHandler) usecases.DBRepository {
	log.Debug("call NewMysqlRepository")
	myRepo := new(MysqlRepo)
	myRepo.myHandler = myHandler
	return myRepo
}

func (myRepo *MysqlRepo) GetUserInfo(id int64) (usecases.User, error) {
	log.Debug("MysqlRepo:GetUserInfo begin id=", id)
	uM, err := myRepo.myHandler.GetUser(id)
	u := createUcUser(uM)
	log.Debugf("MysqlRepo:GetUserInfo result user=%v, err=%v\n", u, err)
	return u, err

}

func (myRepo *MysqlRepo) AddUser(u *usecases.User) error {
	log.Debug("MysqlRepo: AddUser:begin", u)

	uM := createModelUser(*u)
	err := myRepo.myHandler.InsertUser(&uM)
	if err != nil {
		log.Error("MysqlRepo:AddUser err=", err)
		return err
	}
	u.ID = uM.ID
	log.Debug("MysqlRepo:AddUser:success =", u)
	return nil
}

func (myRepo *MysqlRepo) UpdateUser(u usecases.User) error {
	log.Debug("MysqlRepo:UpdateUser:begin user=", u)
	uM := createModelUser(u)
	err := myRepo.myHandler.UpdateUser(uM)
	if err != nil {
		log.Error("MysqlRepo:UpdateUser:error=", err)
		return err
	}
	log.Debug("MysqlRepo:UpdateUser:successful")
	return nil
}
func (myRepo *MysqlRepo) DeleteUser(id int64) error {
	log.Debug("MysqlRepo:DeleteUser:begin id=", id)

	err := myRepo.myHandler.DeleteUser(id)
	if err != nil {
		log.Error("MysqlRepo:DeleteUser:error=", err)
		return err
	}
	log.Debug("MysqlRepo:DeleteUser:successful")
	return nil
}

func (myRepo *MysqlRepo) CloseDB() {
	if err := myRepo.myHandler.CloseDB(); err == nil {
		log.Debug("MysqlRepo: CloseDB() successful")
	} else {
		log.Errorf("MysqlRepo: CloseDB() err %v\n", err)
	}
}

func (myRepo *MysqlRepo) CreateTable(table string) error {
	log.Debug("MysqlRepo:CreateTable ", table)
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
