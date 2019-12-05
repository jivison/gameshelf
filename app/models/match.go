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

// MatchScores return all the match scores associated with a match
func (match Match) MatchScores() []MatchScore {
	var matchScores []MatchScore
	dbmap.Select(matchScores, "select * from match_scores where \"MatchID\"=$1", match.ID)
	return matchScores
}

// AverageScore calculates the average score
func (match Match) AverageScore() float32 {
	var total float32

	scores := match.MatchScores()

	for _, matchScore := range scores {
		total += matchScore.BaseScore
	}

	return total / float32(len(scores))
}

// CalculateAll re-calculates all the matchscores final scores (because the average changes with every new player)
func (match Match) CalculateAll() {
	for _, matchScore := range match.MatchScores() {
		matchScore.CalculateFinalScore()
	}
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
