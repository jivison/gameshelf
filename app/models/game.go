package models

import "log"

// Game is a board game
type Game struct {
	ID    int
	Title string
	Year  int
	BggID int
}

// FindGame finds a game buy its id
func FindGame(id int) (bool, *Game) {
	obj, err := dbmap.Get(Game{}, id)
	game := obj.(*Game)

	if err != nil {
		log.Print("ERROR FindGame: ")
		log.Println(err)
	}

	return (err == nil), game
}

// CreateGame creates a game
func CreateGame(title string, year, bggID int) (bool, *Game) {
	game := &Game{
		Title: title,
		Year:  year,
		BggID: bggID,
	}
	err := dbmap.Insert(game)

	if err != nil {
		log.Print("ERROR: ")
		log.Println(err)
	}

	return (err == nil), game
}
