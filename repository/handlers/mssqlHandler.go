//Copyright © 2018 Fuf
package handlers

import (
	"database/sql"
	"errors"
	"log"

	"github.com/gitfuf/userserver/repository/models"

	_ "github.com/denisenkom/go-mssqldb"
)

type MsHandler struct {
	conn *sql.DB
}

func (msH *MsHandler) InsertUser(u *models.User) error {
	log.Printf("MsHandler: InsertUser begin age=%d, first=%s, last=%s, email=%s\n", u.Age, u.FirstName, u.LastName, u.Email)
	sqlStatement := `
	INSERT INTO users (age, first_name, last_name, email)
	OUTPUT Inserted.id
	VALUES (?, ?, ?, ?)`

	err := msH.conn.QueryRow(sqlStatement, u.Age, u.FirstName, u.LastName, u.Email).Scan(&u.ID)
	if err != nil {
		log.Println("MsHandler: InsertUser err=", err)
		return err
	}
	log.Println("MsHandler: InsertUser successful id=", u.ID)
	return nil
}

func (msH *MsHandler) GetUser(id int64) (models.User, error) {
	log.Println("MsHandler:GetUser begin id=", id)
	var u models.User
	sqlStatement := `SELECT * FROM users WHERE id=?`
	row := msH.conn.QueryRow(sqlStatement, id)
	err := row.Scan(&u.ID, &u.Age, &u.FirstName, &u.LastName, &u.Email)
	if err != nil {
		log.Println("MsHandler:GetUser err=", err)
		switch err {
		case sql.ErrNoRows:
			return u, errors.New("haven't found")
		default:
			return u, err
		}
	}

	log.Println("MsHandler:GetUser success=", u)
	return u, nil
}

func (msH *MsHandler) UpdateUser(u models.User) error {
	log.Printf("MsHandler:UpdateUser begin user=%v \n", u)
	sqlStatement := `
	UPDATE users SET first_name = ?, last_name = ?, email = ?, age = ?
	WHERE id = ?`
	_, err := msH.conn.Exec(sqlStatement, u.FirstName, u.LastName, u.Email, u.Age, u.ID)
	if err != nil {
		log.Println("MsHandler:UpdateUser err=", err)
		return err
	}
	log.Println("MsHandler:UpdateUser successful")
	return nil
}

func (msH *MsHandler) DeleteUser(id int64) error {
	log.Println("MsHandler:DeleteUser begin id=", id)
	sqlStatement := `
	DELETE FROM users
	WHERE id = ?`
	res, err := msH.conn.Exec(sqlStatement, id)
	if err != nil {
		log.Println("MsHandler:DeleteUser err=", err)
		return err
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return errors.New("no such ID")
	}
	log.Println("MsHandler:DeleteUser successful")
	return nil
}

func NewMssqlHandler(connString string) (*MsHandler, error) {
	log.Println("Handlers:NewMssqlHandler:connStr=", connString)
	conn, err := sql.Open("mssql", connString) //only check params
	if err != nil {
		log.Println("mssql err=", err)
		return nil, err
	}
	//check really connect to db
	err = conn.Ping()
	if err != nil {
		log.Println("mssql err=", err)
		conn.Close()
		return nil, err
	}
	MsHandler := new(MsHandler)
	MsHandler.conn = conn
	log.Println("Connect to mssql")
	return MsHandler, nil
}

func (msH *MsHandler) CloseDB() error {
	return msH.conn.Close()
}

//query for create user table
const msTableUserCreateQuery = `
IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='users' AND xtype='U')
	CREATE TABLE users 
	(id INT IDENTITY(1,1) NOT NULL PRIMARY KEY, 
	age INT, 
	first_name NVARCHAR(50), 
	last_name NVARCHAR(50), 
	email NVARCHAR(100) NOT NULL
)`

func (msH *MsHandler) ClearUserTable() error {
	//pgH.conn.Exec("DELETE FROM users")
	_, err := msH.conn.Exec("TRUNCATE TABLE users")
	return err
}

func (msH *MsHandler) CreateUserTable() error {
	_, err := msH.conn.Exec(msTableUserCreateQuery)
	log.Println("MsHandler:CreateUserTable err=", err)
	return err
}
