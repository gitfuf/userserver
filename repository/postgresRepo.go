package repository

import (
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/gitfuf/userserver/repository/handlers"
	"github.com/gitfuf/userserver/usecases"
)

type PostgresRepo struct {
	usecases.DBRepository
	pgHandler *handlers.PostgresHandler
}

func NewPostgresRepository(pgHandler *handlers.PostgresHandler) usecases.DBRepository {
	log.Debug("call NewPostgresRepository")
	pgRepo := new(PostgresRepo)
	pgRepo.pgHandler = pgHandler
	return pgRepo
}

func (pgRepo *PostgresRepo) GetUserInfo(id int64) (usecases.User, error) {
	log.Debug("PostgresRepo:GetUserInfo begin id=", id)
	uM, err := pgRepo.pgHandler.GetUser(id)
	u := createUcUser(uM)
	log.Debugf("PostgresRepo:GetUserInfo result user=%v, err=%v\n", u, err)
	return u, err

}

func (pgRepo *PostgresRepo) AddUser(u *usecases.User) error {
	log.Debug("PostgresRepo: AddUser:begin", u)

	uM := createModelUser(*u)
	err := pgRepo.pgHandler.InsertUser(&uM)
	if err != nil {
		log.Error("PostgresRepo:AddUser err=", err)
		return err
	}
	u.ID = uM.ID
	log.Debug("PostgresRepo:AddUser:success =", u)
	return nil
}

func (pgRepo *PostgresRepo) UpdateUser(u usecases.User) error {
	log.Debug("PostgresRepo:UpdateUser:begin user=", u)
	uM := createModelUser(u)
	err := pgRepo.pgHandler.UpdateUser(uM)
	if err != nil {
		log.Error("PostgresRepo:UpdateUser:error=", err)
		return err
	}
	log.Debug("PostgresRepo:UpdateUser:successful")
	return nil
}
func (pgRepo *PostgresRepo) DeleteUser(id int64) error {
	log.Debug("PostgresRepo:DeleteUser:begin id=", id)

	err := pgRepo.pgHandler.DeleteUser(id)
	if err != nil {
		log.Error("PostgresRepo:DeleteUser:error=", err)
		return err
	}
	log.Debug("PostgresRepo:DeleteUser:successful")
	return nil
}

func (pgRepo *PostgresRepo) CloseDB() {
	if err := pgRepo.pgHandler.CloseDB(); err == nil {
		log.Debug("PostgresRepo: CloseDB() successful")
	} else {
		log.Errorf("PostgresRepo: CloseDB() err %v\n", err)
	}

}

func (pgRepo *PostgresRepo) CreateTable(table string) error {
	log.Debug("PostgresRepo:CreateTable ", table)
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
