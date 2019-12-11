package models

import "log"

// GroupGame is a join table between groups and games
type GroupGame struct {
	ID          int
	GameID      int
	GroupID     int
	TimesPlayed int
}

// CreateGroupGame creates a groupGame in the database
func CreateGroupGame(groupID, gameID int) (bool, *GroupGame) {
	groupGame := &GroupGame{
		GroupID: groupID,
		GameID:  gameID,
	}

	err := dbmap.Insert(groupGame)

	if err != nil {
		log.Print("CreateGroupGame ERROR:")
		log.Print(err)
	}

	return (err == nil), groupGame
}
