package repository

import (
	"github.com/gitfuf/userserver/repository/models"
	"github.com/gitfuf/userserver/usecases"

	"log"
)

func createModelUser(u usecases.User) models.User {
	mu := models.User{
		ID:        u.ID,
		Age:       u.Age,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}
	log.Printf("createModelUser %v from %v\n", mu, u)
	return mu
}

func createUcUser(u models.User) usecases.User {
	uc := usecases.User{
		ID:        u.ID,
		Age:       u.Age,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}
	log.Printf("createUcUser %v from %v\n", uc, u)
	return uc
}
