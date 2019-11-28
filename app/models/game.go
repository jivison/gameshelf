package models

import "log"

// Game is a board game
type Game struct {
	ID       int
	Title    string
	Year     int
	BggID    int
	Username string `db:"user_name"`
	User     *User  `db:"-"`
}

// FindGame finds a game buy its id
func FindGame(id int) (bool, *Game) {
	obj, err := dbmap.Get(Game{}, id)
	game := obj.(*Game)

	if err != nil {
		log.Print("ERROR FindGame: ")
		log.Println(err)
	}

	_, game.User = FindUser(game.Username)

	return (err == nil), game
}

func FindGameByTitle(title, username string, storageVar *[]Game) {
	dbmap.Select(storageVar, "select * from games where title=$1 and username=$1", title, username)
}

// CreateGame creates a game
func CreateGame(title string, year, bggID int, username string) (bool, *Game) {
	game := &Game{
		Title:    title,
		Year:     year,
		BggID:    bggID,
		Username: username,
	}
	err := dbmap.Insert(game)

	if err != nil {
		log.Print("ERROR: ")
		log.Println(err)
	}

	return (err == nil), game
}
