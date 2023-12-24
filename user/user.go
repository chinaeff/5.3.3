package user

type User struct {
	Email    string
	Password string
}

type UserService interface {
	Register(email, password string) error
	Authenticate(email, password string) (bool, error)
}

func (u *User) Register(email, password string) User {
	return User{Email: email, Password: password}
}
func (u *User) Authenticate(email, password string) bool {
	if u.Email == email && u.Password == password {
		return true
	}
	return false
}
