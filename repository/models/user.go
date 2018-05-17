//Copyright © 2018 Fuf
package models

type User struct {
	ID        int64  `json:"id"`
	Age       int    `json:"age"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
