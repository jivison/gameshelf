package models

import "log"

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
	Game              Game  `db:"-"`
	Match             Match `db:"-"`
}

// CalculateFinalScore sets a match score's final score
func (matchScore MatchScore) CalculateFinalScore() float32 {
	if matchScore.IsWinner {
		matchScore.FinalScore = (matchScore.BaseScore / matchScore.Match.AverageScore()) * 1.10
	}
	matchScore.FinalScore = matchScore.BaseScore / matchScore.Match.AverageScore()
	log.Printf("Average Score: %f", matchScore.Match.AverageScore())
	dbmap.Update(&matchScore)
	return matchScore.FinalScore
}

// ComplexityRating returns the complexity rating for the associated game
func (matchScore MatchScore) ComplexityRating() float32 {
	ok, game := FindGame(matchScore.GameID)
	if ok {
		return game.ComplexityRating
	}
	return 1.0
}

// GetGame sets the Game field on a matchScore to the one that corresponds with its GameID
func (matchScore MatchScore) GetGame() Game {
	ok, game := FindGame(matchScore.GameID)
	if ok {
		matchScore.Game = *game
		return *game
	}
	return Game{}
}

// GetMatch sets the Match firld on a matchScore to the one that correspondes with its MatchID
func (matchScore MatchScore) GetMatch() Match {
	ok, match := FindMatch(matchScore.MatchID)
	if ok {
		matchScore.Match = *match
		return *match
	}
	return Match{}
}

func CreateMatchScore(match Match, game Game, playerUserName string, baseScore float32, isWinner bool) (bool, *MatchScore) {
	ok, player := FindUser(playerUserName)

	if !ok {
		return false, nil
	}

	matchScore := &MatchScore{
		MatchID:           match.ID,
		GameID:            game.ID,
		PlayerUserName:    playerUserName,
		PlayerDisplayName: player.FirstName,
		BaseScore:         baseScore,
		IsWinner:          isWinner,
	}

	matchScore.Match = match
	matchScore.Game = game

	err := dbmap.Insert(matchScore)

	if err != nil {
		log.Print("ERROR: ")
		log.Println(err)
	}

	match.CalculateAll()

	return (err == nil), matchScore
}
