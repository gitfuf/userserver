//Copyright © 2018 Fuf
package models

import "database/sql"

type User struct {
	ID        int64          `json:"id"`
	Age       sql.NullInt64  `json:"age"`
	FirstName sql.NullString `json:"first_name"`
	LastName  sql.NullString `json:"last_name"`
	Email     string         `json:"email"`
}
