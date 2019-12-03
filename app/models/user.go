package models

import (
	"fmt"
	"log"
)

// User is a user
type User struct {
	Username       string `db:",primarykey"`
	Password       string `db:"-"`
	HashedPassword []byte
}

func (u User) String() string {
	return fmt.Sprintf("{ USERNAME: %s | HASHEDPASSWORD: %s }", u.Username, u.HashedPassword)
}

func (u User) Games() []Game {
	var games []Game
	dbmap.Select(&games, "select * from games where user_name=$1", u.Username)
	return games
}

// FindUser finds a user by its username
func FindUser(username string) (bool, *User) {
	var users []User
	_, err := dbmap.Select(&users, "select * from users where \"Username\"=$1", username)

	if len(users) < 1 {
		log.Printf("INFO: FindUser couldn't find a user with username %s", username)
		return false, &User{}
	}

	user := users[0]

	if err != nil {
		log.Print("ERROR FindUser: ")
		log.Println(err)
	}

	return (err == nil), &user
}

// InsertUser inserts an already 'built' user
func InsertUser(user User) bool {
	err := dbmap.Insert(&user)

	if err != nil {
		log.Print("ERROR: ")
		log.Println(err)
	}

	return err == nil

}
