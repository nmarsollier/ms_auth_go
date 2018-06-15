package user

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"golang.org/x/crypto/bcrypt"
)

// User data structure
type User struct {
	ID          objectid.ObjectID `bson:"_id"`
	Name        string            `bson:"name"`
	Login       string            `bson:"login"`
	Password    string            `bson:"password"`
	Permissions []string          `bson:"permissions"`
	Enabled     bool              `bson:"enabled"`
	Created     time.Time         `bson:"created"`
	Updated     time.Time         `bson:"updated"`
}

func newUser() *User {
	return &User{
		ID:          objectid.New(),
		Enabled:     true,
		Created:     time.Now(),
		Updated:     time.Now(),
		Permissions: []string{"user"},
	}
}

func (e *User) setPasswordText(pwd string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return ErrPassword
	}

	e.Password = string(hash)
	return nil
}

func (e *User) validatePassword(plainPwd string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(e.Password), []byte(plainPwd)); err != nil {
		return ErrPassword
	}
	return nil
}

// Granted verifica si el usuario tiene el permiso indicado
func (e *User) granted(permission string) bool {
	for _, p := range e.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// Grant le otorga el permiso indicado al usuario
func (e *User) grant(permission string) {
	if !e.granted(permission) {
		e.Permissions = append(e.Permissions, permission)
	}
}

// Revoke le revoca el permiso indicado al usuario
func (e *User) revoke(permission string) {
	if e.granted(permission) {
		var newPermissions []string
		for _, p := range e.Permissions {
			if p != permission {
				newPermissions = append(newPermissions, p)
			}
		}
		e.Permissions = newPermissions
	}
}
