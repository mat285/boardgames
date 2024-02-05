package v1alpha1

import "github.com/blend/go-sdk/uuid"

type User struct {
	Object
	Username string
}

func NewUser(id uuid.UUID, username string) User {
	return User{
		Object:   NewObject(id),
		Username: username,
	}
}

func (u User) GetUsername() string {
	return u.Username
}
