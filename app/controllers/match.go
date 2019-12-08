package controllers

import (
	"fmt"
	"gameshelf/app/models"
	"sort"
	"time"

	"github.com/revel/revel"
)

// Match controller handles creation, updation and deletion of matches
type Match struct {
	*revel.Controller
}

// New displays the form to create a new match
func (c Match) New(gameid int) revel.Result {
	ok, game := models.FindGame(gameid)
	if ok {
		return c.Render(game)
	}
	return c.RenderText(fmt.Sprintf("Couldn't find a game with that id! (%d)", gameid))
}

// Create creates a match in the database
func (c Match) Create(gameid int, datePlayed time.Time) revel.Result {
	username, _ := c.Session.Get("user")

	c.Validation.Required(datePlayed).Message("Must have a date played")

	c.Validation.Required(username != nil).Message("You must be signed in to create a match")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Match.New)
	}

	if username != nil {
		if ok, match := models.CreateMatch(gameid, username.(string), datePlayed); ok {
			return c.Redirect(Match.Show, match.GameID, match.ID)
		}
	}

	c.FlashParams()
	return c.Redirect(Match.New)
}

// Show displays a match
func (c Match) Show(id int) revel.Result {
	ok, match := models.FindMatch(id)

	if ok {
		scores := match.MatchScores()
		sort.SliceStable(scores, func(i, j int) bool {
			if (scores[i].IsWinner && !scores[j].IsWinner) || (!scores[i].IsWinner && scores[j].IsWinner) {
				return scores[i].IsWinner
			}
			return scores[i].FinalScore > scores[j].FinalScore
		})

		return c.Render(scores, match)
	}

	return c.RenderText(fmt.Sprintf("Couldn't find a match with that id! (%d)", id))
}

// AddScore adds a match score to a match
func (c Match) AddScore(gameid, id int, playerUserName string, baseScore float32, isWinner bool) revel.Result {
	ok, game := models.FindGame(gameid)
	if ok {
		ok, match := models.FindMatch(id)
		if ok {
			players := match.Players()

			if players[playerUserName] {
				c.Flash.Error("That player already has a score")
				return c.Redirect(Match.Show, match.GameID, match.ID)
			}

			if ok, _ := models.CreateMatchScore(match, game, playerUserName, baseScore, isWinner); !ok {
				c.Flash.Error("Couldn't find a user with that username")
				c.FlashParams()
			}
			return c.Redirect(Match.Show, match.GameID, match.ID)
		}

		return c.RenderText(fmt.Sprintf("Couldn't find a match with that id! (%d)", id))
	}

	return c.RenderText(fmt.Sprintf("Couldn't find a game with that id! (%d)", gameid))
}

// RemoveScore removes a match score from a match
func (c Match) RemoveScore(gameid, mid, id int) revel.Result {
	ok, matchScore := models.FindMatchScore(id)
	if ok {
		matchScore.Delete()

		_, match := models.FindMatch(mid)
		match.CalculateAll()
	}

	return c.Redirect(Match.Show, gameid, mid)
}

// Delete delets a match from the database
func (c Match) Delete(id int) revel.Result {
	ok, match := models.FindMatch(id)
	if ok {
		match.Delete()
		return c.Redirect(Game.Show, match.GameID)
	}
	return c.RenderText(fmt.Sprintf("Couldn't find a match with that id! (%d)", id))
}
