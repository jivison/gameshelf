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
	Games          []Game `db:"-"`
}

func (u User) String() string {
	return fmt.Sprintf("{ USERNAME: %s | PASSWORD: [HIDDEN] | HASHEDPASSWORD: %s }", u.Username, u.HashedPassword)
}

// AddGame adds a game to the user structs games
func (u User) AddGame(game Game) {
	u.Games = append(u.Games, game)
}

// FindUser finds a user by its username
func FindUser(username string) (bool, *User) {
	// obj, err := dbmap.Get(User{}, username)
	// user := obj.(*User)

	var users []User
	_, err := dbmap.Select(&users, "select * from users where \"Username\"=$1", username)

	user := users[0]

	if err != nil {
		log.Print("ERROR FindUser: ")
		log.Println(err)
	} else {
		var games []Game
		_, err = dbmap.Select(games, "select * from users where user_name=$1", user.Username)
		user.Games = games
	}

	return (err != nil), &user
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
