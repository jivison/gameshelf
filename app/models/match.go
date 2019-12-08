package models

import (
	"log"
	"math"
	"time"
)

// Match is a game of a board game
type Match struct {
	ID           int
	GameID       int
	DatePlayed   time.Time
	HostUserName string
}

// Delete deletes a match from the database
func (m Match) Delete() bool {
	_, err := dbmap.Delete(&m)
	return (err != nil)
}

// MatchScores return all the match scores associated with a match
func (m Match) MatchScores() []MatchScore {
	var matchScores []MatchScore
	dbmap.Select(&matchScores, "select * from match_scores where \"MatchID\"=$1", m.ID)
	return matchScores
}

// Players returns a slice of all the players in a game
func (m Match) Players() map[string]bool {
	scores := m.MatchScores()

	players := make(map[string]bool)

	for _, score := range scores {
		players[score.PlayerUserName] = true
	}

	return players
}

// AverageScore calculates the average score
func (m Match) AverageScore() float32 {
	var total float32

	scores := m.MatchScores()

	if len(scores) == 1 {
		return float32(math.Inf(1))
	}

	for _, matchScore := range scores {
		total += matchScore.BaseScore
	}

	result := total / float32(len(scores))

	return result
}

// CalculateAll re-calculates all the matchscores final scores (because the average changes with every new player)
func (m Match) CalculateAll() {
	for _, matchScore := range m.MatchScores() {
		matchScore.CalculateFinalScore(m)
	}
}

// Game returns the game associated with a game
func (m Match) Game() *Game {
	_, game := FindGame(m.GameID)
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
