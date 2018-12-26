package root

type User struct {
	Id       string `json:"_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserService interface {
	CreateUser(u *User) (string, error)
	GetAllUsers() ([]*User, error)
	GetUserById(id string) (*User, error)
	UpdateUser(u *User) error
	DeleteUserById(id string) error
}
