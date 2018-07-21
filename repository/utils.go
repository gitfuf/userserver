package repository

import (
	"github.com/gitfuf/userserver/repository/models"
	"github.com/gitfuf/userserver/usecases"
	log "github.com/sirupsen/logrus"

	"database/sql"
)

func createModelUser(u usecases.User) models.User {
	//Because of NULL values need to do convertion
	var (
		age         sql.NullInt64
		first, last sql.NullString
	)
	age.Int64 = u.Age
	if u.Age > 0 {
		age.Valid = true
	}
	first.String = u.FirstName
	if u.FirstName != "" {
		first.Valid = true
	}
	last.String = u.LastName
	if u.LastName != "" {
		last.Valid = true
	}

	mu := models.User{
		ID:        u.ID,
		Age:       age,
		FirstName: first,
		LastName:  last,
		Email:     u.Email,
	}

	log.Debugf("createModelUser %v from %v\n", mu, u)
	return mu
}

func createUcUser(u models.User) usecases.User {
	uc := usecases.User{
		ID:        u.ID,
		Age:       u.Age.Int64,
		FirstName: u.FirstName.String,
		LastName:  u.LastName.String,
		Email:     u.Email,
	}
	log.Debugf("createUcUser %v from %v\n", uc, u)
	return uc
}
