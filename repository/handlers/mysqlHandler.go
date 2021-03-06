//Copyright © 2018 Fuf
//mysql infrastructure
package handlers

import (
	"database/sql"
	"errors"

	"github.com/gitfuf/userserver/repository/models"
	log "github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlHandler struct {
	conn *sql.DB
}

func NewMysqlHandler(connString string) (*MysqlHandler, error) {
	log.Debug("call NewMysqlHandler str:", connString)
	conn, err := sql.Open("mysql", connString) //only check params
	if err != nil {
		log.Error("err=", err)
		return nil, err
	}
	//check really connect to db
	err = conn.Ping()
	if err != nil {
		conn.Close()
		log.Error("db ping err=", err)
		return nil, err
	}

	mysqlHandler := new(MysqlHandler)
	mysqlHandler.conn = conn
	log.Infof("NewMysqlHandler: connect to db")
	return mysqlHandler, nil
}

func (myH *MysqlHandler) InsertUser(u *models.User) error {
	log.Debug("MysqlHandler: InsertUser begin age=%d, first=%s, last=%s, email=%s\n", u.Age, u.FirstName, u.LastName, u.Email)

	sqlStatement := `
	INSERT INTO users (age, first_name, last_name, email)
	VALUES (?, ?, ?, ?)`

	stmt, err := myH.conn.Prepare(sqlStatement)
	if err != nil {
		log.Error("MysqlHandler: InsertUser prepare statement err=", err)
		return err
	}
	res, err := stmt.Exec(u.Age, u.FirstName, u.LastName, u.Email)
	if err != nil {
		log.Error("MysqlHandler: InsertUser exec statement err=", err)
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Error("MysqlHandler: InsertUser LastInsertId err=", err)

	} else {
		u.ID = id
	}

	log.Debug("MysqlHandler: InsertUser successful id=", u.ID)
	return nil
}

func (myH *MysqlHandler) GetUser(id int64) (models.User, error) {
	log.Debug("MysqlHandler:GetUser begin id=", id)
	var u models.User
	sqlStatement := `SELECT id, age, first_name, last_name, email FROM users WHERE id=?;`
	row := myH.conn.QueryRow(sqlStatement, id)
	err := row.Scan(&u.ID, &u.Age, &u.FirstName, &u.LastName, &u.Email)
	if err != nil {
		log.Error("MysqlHandler:GetUser err=", err)
		switch err {
		case sql.ErrNoRows:
			return u, errors.New("haven't found")
		default:
			return u, err
		}
	}

	log.Debug("MysqlHandler:GetUser success=", u)
	return u, nil
}

func (myH *MysqlHandler) UpdateUser(u models.User) error {
	log.Debugf("MysqlHandler:UpdateUser begin user=%v \n", u)
	sqlStatement := `
	UPDATE users SET first_name = ?, last_name = ?, email = ?, age = ?
	WHERE id = ?;`
	_, err := myH.conn.Exec(sqlStatement, u.FirstName, u.LastName, u.Email, u.Age, u.ID)
	if err != nil {
		log.Error("MysqlHandler:UpdateUser err=", err)
		return err
	}
	log.Debug("MysqlHandler:UpdateUser successful")
	return nil
}

func (myH *MysqlHandler) DeleteUser(id int64) error {
	log.Debug("MysqlHandler:DeleteUser begin id=", id)
	sqlStatement := `
	DELETE FROM users
	WHERE id = ?;`
	res, err := myH.conn.Exec(sqlStatement, id)
	if err != nil {
		log.Error("MysqlHandler:DeleteUser err=", err)
		return err
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return errors.New("no such ID")
	}
	log.Debug("MysqlHandler:DeleteUser successful")

	return nil
}

func (myH *MysqlHandler) CloseDB() error {
	return myH.conn.Close()
}

//query for create user table
const myTableUserCreateQuery = `CREATE TABLE IF NOT EXISTS users (
        id INT(10) NOT NULL AUTO_INCREMENT,
		age INT NULL DEFAULT NULL,
        first_name VARCHAR(64) NULL DEFAULT NULL,
        last_name VARCHAR(64) NULL DEFAULT NULL,
        email VARCHAR(64) NOT NULL,
        PRIMARY KEY (id)
    );`

func (myH *MysqlHandler) ClearUserTable() error {
	_, err := myH.conn.Exec("TRUNCATE users ")
	//_, err = myH.conn.Exec("ALTER TABLE users AUTO_INCREMENT = 1")
	return err
}

func (myH *MysqlHandler) CreateUserTable() error {
	_, err := myH.conn.Exec(myTableUserCreateQuery)
	log.Debug("MysqlHandler:CreateUserTable err=", err)
	return err
}
