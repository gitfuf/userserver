package usecases

type DBRepository interface {
	UserRepository
	RoleRepository
	CreateTable(tableName string) error
	ClearTable(tableName string) error
}

type UserRepository interface {
	AddUser(u *User) error
	GetUserInfo(id int64) (User, error)
	UpdateUser(u User) error
	DeleteUser(id int64) error
	//SelectAllEmails() ([]string, error)
}

type RoleRepository interface {
	AddRole() error
}
type User struct {
	ID        int64  `json:"id"`
	Age       int    `json:"age"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
