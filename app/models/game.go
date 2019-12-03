package models

import (
	"fmt"
	"log"
)

// Game is a board game
type Game struct {
	ID       int
	Title    string
	Year     int
	BggID    int
	ImgURL   string
	Username string `db:"user_name"`
	User     *User  `db:"-"`
}

func (game Game) Refresh() {

}

func (game Game) String() string {
	return fmt.Sprintf("{ ID: %d | TITLE: %s | YEAR: %d | BGGID: %d | USERNAME: %s | IMGURL: %s }", game.ID, game.Title, game.Year, game.BggID, game.Username, game.ImgURL)
}

// Delete deletes a game from the database
func (game Game) Delete() bool {
	_, err := dbmap.Delete(&game)
	return (err != nil)
}

// Update updates the database with any changes to a Game
func (game Game) Update() error {
	_, err := dbmap.Update(&game)
	return err
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

// FindGameByTitle finds a game by its title
func FindGameByTitle(title, username string, storageVar *[]Game) {
	dbmap.Select(storageVar, "select * from games where title=$1 and username=$1", title, username)
}

// CreateGame creates a game
func CreateGame(title string, year, bggID int, username, imgURL string) (bool, *Game) {
	game := &Game{
		Title:    title,
		Year:     year,
		BggID:    bggID,
		Username: username,
		ImgURL:   imgURL,
	}

	err := dbmap.Insert(game)

	if err != nil {
		log.Print("ERROR: ")
		log.Println(err)
	}

	return (err == nil), game
}
