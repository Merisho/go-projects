package auth

import "errors"

func NewAuthenticator() *Authenticator {
	return &Authenticator{
		users: make(map[string]*User),
	}
}

type Authenticator struct {
	users map[string]*User
}

func (a *Authenticator) Auth(username, password string) error {
	user := a.getUser(username)
	if user == nil {
		a.createUser(username, password)
		return nil
	}

	if user.Password == password {
		return nil
	}

	return errors.New("wrong password")
}

func (a *Authenticator) getUser(username string) *User {
	return a.users[username]
}

func (a *Authenticator) createUser(username, password string) {
	a.users[username] = &User{
		Username: username,
		Password: password,
	}
}
