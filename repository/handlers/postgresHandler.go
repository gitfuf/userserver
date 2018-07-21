//Copyright © 2018 Fuf
//postgres infrastructure
package handlers

import (
	"database/sql"
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/gitfuf/userserver/repository/models"

	_ "github.com/lib/pq"
)

type PostgresHandler struct {
	conn *sql.DB
}

func NewPostgresHandler(connString string) (*PostgresHandler, error) {
	log.Debug("call NewPostgresHandler str:", connString)
	conn, err := sql.Open("postgres", connString) //only check params
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
	postgresHandler := new(PostgresHandler)
	postgresHandler.conn = conn
	log.Info("NewPostgresHandler: connect to db")
	return postgresHandler, nil
}

func (pgH *PostgresHandler) InsertUser(u *models.User) error {
	log.Debugf("PostgresHandler: InsertUser begin age=%d, first=%s, last=%s, email=%s\n", u.Age, u.FirstName, u.LastName, u.Email)
	sqlStatement := `
	INSERT INTO users (age, first_name, last_name, email)
	VALUES ($1, $2, $3, $4)
	RETURNING id`

	err := pgH.conn.QueryRow(sqlStatement, u.Age, u.FirstName, u.LastName, u.Email).Scan(&u.ID)
	if err != nil {
		log.Error("PostgresHandler: InsertUser err=", err)
		return err
	}
	log.Debug("PostgresHandler: InsertUser successful id=", u.ID)
	return nil
}

func (pgH *PostgresHandler) GetUser(id int64) (models.User, error) {
	log.Debug("PostgresHandler:GetUser begin id=", id)
	var u models.User
	sqlStatement := `SELECT * FROM users WHERE id=$1;`
	row := pgH.conn.QueryRow(sqlStatement, id)
	err := row.Scan(&u.ID, &u.Age, &u.FirstName, &u.LastName, &u.Email)
	if err != nil {
		log.Error("PostgresHandler:GetUser err=", err)
		switch err {
		case sql.ErrNoRows:
			return u, errors.New("haven't found")
		default:
			return u, err
		}
	}

	log.Debug("PostgresHandler:GetUser success=", u)
	return u, nil
}

func (pgH *PostgresHandler) UpdateUser(u models.User) error {
	log.Debugf("PostgresHandler:UpdateUser begin user=%v \n", u)
	sqlStatement := `
	UPDATE users SET first_name = $2, last_name = $3, email = $4, age = $5
	WHERE id = $1;`
	_, err := pgH.conn.Exec(sqlStatement, u.ID, u.FirstName, u.LastName, u.Email, u.Age)
	if err != nil {
		log.Error("PostgresHandler:UpdateUser err=", err)
		return err
	}
	log.Debug("PostgresHandler:UpdateUser successful")
	return nil
}

func (pgH *PostgresHandler) DeleteUser(id int64) error {
	log.Debug("PostgresHandler:DeleteUser begin id=", id)
	sqlStatement := `
	DELETE FROM users
	WHERE id = $1;`
	res, err := pgH.conn.Exec(sqlStatement, id)
	if err != nil {
		log.Error("PostgresHandler:DeleteUser err=", err)
		return err
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return errors.New("no such ID")
	}
	log.Debug("PostgresHandler:DeleteUser successful")

	return nil
}

/*
func (pgH *PostgresHandler) SelectAllEmails() ([]string, error) {
	rows, err := dbw.db.Query("SELECT email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	emailList := make([]string, 0, 20)
	for rows.Next() {
		var email string
		err = rows.Scan(&email)
		if err != nil {
			return nil, err
		}
		emailList = append(emailList, email)
	}
	//check that there were no REAL error
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return emailList, nil
}
*/

func (pgH *PostgresHandler) CloseDB() error {
	return pgH.conn.Close()
}

//query for create user table
const pgTableUserCreateQuery = `CREATE TABLE IF NOT EXISTS users
( 	id SERIAL PRIMARY KEY,
	age INT,
	first_name TEXT,
	last_name TEXT,
	email TEXT NOT NULL
)`

func (pgH *PostgresHandler) ClearUserTable() error {
	//pgH.conn.Exec("DELETE FROM users")
	_, err := pgH.conn.Exec("TRUNCATE users RESTART IDENTITY")
	return err
}

func (pgH *PostgresHandler) CreateUserTable() error {
	_, err := pgH.conn.Exec(pgTableUserCreateQuery)
	log.Debug("PostgresHandler:CreateUserTable err=", err)
	return err
}
