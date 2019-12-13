package models

import (
	"log"
	"math"
)

// MatchScore represents a single players score in a match
type MatchScore struct {
	ID                int
	MatchID           int
	GameID            int
	PlayerUserName    string
	PlayerDisplayName string // So we don't have to query the database another time to get the players first name
	BaseScore         float32
	IsWinner          bool
	FinalScore        float32
}

// Delete deletes a matchScore from the database
func (ms MatchScore) Delete() bool {
	_, err := dbmap.Delete(&ms)
	return (err != nil)
}

// CalculateFinalScore sets a match score's final score
func (ms MatchScore) CalculateFinalScore(match Match) float32 {
	average := match.AverageScore()

	game := match.Game()

	if average == float32(math.Inf(1)) {
		ms.FinalScore = 100
	} else if ms.IsWinner {
		ms.FinalScore = (ms.BaseScore / average) * 110
	} else {
		ms.FinalScore = ms.BaseScore / average * 100
	}

	ms.FinalScore = (game.ComplexityRating / 5) * ms.FinalScore

	dbmap.Update(&ms)
	return ms.FinalScore
}

// ComplexityRating returns the complexity rating for the associated game
func (ms MatchScore) ComplexityRating() float32 {
	ok, game := FindGame(ms.GameID)
	if ok {
		return game.ComplexityRating
	}
	return 1.0
}

// CreateMatchScore creates a match score in the database
func CreateMatchScore(match *Match, gameID int, playerUserName string, baseScore float32, isWinner bool) (bool, *MatchScore) {
	ok, player := FindUser(playerUserName)

	if !ok {
		return false, nil
	}

	matchScore := &MatchScore{
		MatchID:           match.ID,
		GameID:            gameID,
		PlayerUserName:    player.Username,
		PlayerDisplayName: player.FirstName,
		BaseScore:         baseScore,
		IsWinner:          isWinner,
	}

	err := dbmap.Insert(matchScore)

	if err != nil {
		log.Print("ERROR: ")
		log.Println(err)
	}

	match.CalculateAll()

	return (err == nil), matchScore
}

// FindMatchScore finds a match score by its ID
func FindMatchScore(id int) (bool, *MatchScore) {
	obj, err := dbmap.Get(MatchScore{}, id)

	var matchScore *MatchScore

	if err != nil {
		log.Print("ERROR FindMatchScore: ")
		log.Println(err)
	} else {
		matchScore = obj.(*MatchScore)
	}

	return (err == nil), matchScore
}
