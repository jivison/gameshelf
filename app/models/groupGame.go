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

// FindGroupGameFromIDs returns a group game from a game and a group ID
func FindGroupGameFromIDs(gameID, groupID int) GroupGame {
	var groupGames []GroupGame
	dbmap.Select(&groupGames, "select * from group_games where \"GameID\"=:gameid and \"GroupID\"=:groupid", map[string]interface{}{
		"gameid":  gameID,
		"groupid": groupID,
	})
	if len(groupGames) > 0 {
		return groupGames[0]
	}
	return GroupGame{}
}

// FindGroupGame finds a group game from its ID
func FindGroupGame(id int) (bool, *GroupGame) {
	var groupGames []GroupGame

	dbmap.Select(&groupGames, findQstring("group_games"), id)

	if len(groupGames) > 0 {
		return true, &groupGames[0]
	}
	return false, &GroupGame{}

}

// GroupIDFromGroupGameID returns a groups id from a group game ID
func GroupIDFromGroupGameID(groupGameID int) int {
	ok, gg := FindGroupGame(groupGameID)
	if ok {
		return gg.GroupID
	}
	return 0
}
