package models

import (
	"log"
	"time"
)

// Match is a game of a board game
type Match struct {
	ID           int
	GameID       int
	DatePlayed   time.Time
	HostUserName string
}

// Game returns the game associated with a game
func (match Match) Game() Game {
	var game Game
	dbmap.Select(game, "select * from games where \"ID\"=$1", match.Game)
	return game
}

// CreateMatch creates a match in the database
func CreateMatch(gameID int, hostUserName string, datePlayed time.Time) (bool, *Match) {
	match := &Match{
		GameID:       gameID,
		HostUserName: hostUserName,
		DatePlayed:   datePlayed,
	}

	err := dbmap.Insert(match)

	if err != nil {
		log.Print("ERROR: ")
		log.Println(err)
	}

	return (err == nil), match
}

// FindMatch finds a match from its id
func FindMatch(id int) (bool, *Match) {
	obj, err := dbmap.Get(Match{}, id)

	var match *Match

	if err != nil {
		log.Print("ERROR FindMatch: ")
		log.Println(err)
	} else {
		match = obj.(*Match)
	}

	return (err == nil), match
}
