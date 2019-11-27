package models

import "log"

// User is a user
type User struct {
	Username       string `db:",primarykey"`
	Password       string
	FirstName      string
	HashedPassword []byte
	Games          []Game `db:"-"`
}

// AddGame adds a game to the user structs games
func (u User) AddGame(game Game) {
	u.Games = append(u.Games, game)
}

// FindUser finds a user by its username
func FindUser(username string) (bool, *User) {
	obj, err := dbmap.Get(User{}, username)
	user := obj.(*User)

	if err != nil {
		log.Print("ERROR FindUser: ")
		log.Println(err)
	} else {
		var games []Game
		_, err = dbmap.Select(games, "select * from users where user_name=$1", user.Username)
		user.Games = games
	}

	return (err == nil), user
}

// FindUserByEmail finds a user by their email
func FindUserByEmail(email string) (bool, *User) {
	var user User
	_, err := dbmap.Select(&user, "select * from users where email=$1", email)
	return err == nil, &user
}

// InsertUser inserts an already 'built' user
func InsertUser(user User) bool {
	err := dbmap.Insert(user)

	if err != nil {
		log.Print("ERROR: ")
		log.Println(err)
	}

	return err == nil

}
