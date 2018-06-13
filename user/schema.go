package user

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/nmarsollier/authgo/tools/lookup"
	"golang.org/x/crypto/bcrypt"
)

// User data structure
type User struct {
	_id      string
	Name     string    `bson:"name"`
	Login    string    `bson:"login"`
	Password string    `bson:"password"`
	Roles    []string  `bson:"roles"`
	Enabled  bool      `bson:"enabled"`
	Created  time.Time `bson:"created"`
	Updated  time.Time `bson:"updated"`
}

func newUser() *User {
	return &User{
		Enabled: true,
		Created: time.Now(),
		Updated: time.Now(),
	}
}

func (e *User) setID(ID objectid.ObjectID) {
	e._id = ID.Hex()
}

// ID obtiene el id de usuario
func (e *User) ID() string {
	return e._id
}

func newUserFromBson(document bson.Document) *User {
	return &User{
		_id:      lookup.ObjectID(document, "_id"),
		Login:    lookup.String(document, "login"),
		Name:     lookup.String(document, "name"),
		Password: lookup.String(document, "password"),
		Enabled:  lookup.Bool(document, "enable"),
		Roles:    lookup.StringArray(document, "roles"),
		Created:  lookup.Time(document, "created"),
		Updated:  lookup.Time(document, "updated"),
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
	err := bcrypt.CompareHashAndPassword([]byte(e.Password), []byte(plainPwd))

	if err != nil {
		return ErrPassword
	}
	return nil
}
