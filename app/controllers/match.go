package controllers

import (
	"fmt"
	"gameshelf/app/models"
	"log"
	"sort"
	"time"

	"github.com/revel/revel"
)

// Match controller handles creation, updation and deletion of matches
type Match struct {
	*revel.Controller
}

// New displays the form to create a new match
func (c Match) New(gid, id int) revel.Result {
	ok, game := models.FindGame(id)
	if ok {
		return c.Render(game, gid)
	}
	return c.RenderText(fmt.Sprintf("Couldn't find a game with that id! (%d)", id))
}

// Create creates a match in the database
func (c Match) Create(gid, id int, datePlayed time.Time) revel.Result {
	username, _ := c.Session.Get("user")

	c.Validation.Required(datePlayed).Message("Must have a date played")

	c.Validation.Required(username != nil).Message("You must be signed in to create a match")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Match.New, gid, id)
	}

	if username != nil {
		if ok, match := models.CreateMatch(id, gid, username.(string), datePlayed); ok {
			return c.Redirect(Match.Show, match.ID)
		}
	}

	c.FlashParams()
	return c.Redirect(Match.New, gid, id)
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

		groupid := models.GroupIDFromGroupGameID(match.GroupGameID)

		return c.Render(scores, match, groupid)
	}

	return c.RenderText(fmt.Sprintf("Couldn't find a match with that id! (%d)", id))
}

// AddScore adds a match score to a match
func (c Match) AddScore(id int, playerUserName string, baseScore float32, isWinner bool) revel.Result {
	ok, match := models.FindMatch(id)
	if ok {
		players := match.Players()

		if players[playerUserName] {
			c.Flash.Error("That player already has a score")
			return c.Redirect(Match.Show, match.GameID, match.ID)
		}

		if ok, _ := models.CreateMatchScore(match, match.GameID, playerUserName, baseScore, isWinner); !ok {
			c.Flash.Error("Couldn't find a user with that username")
			c.FlashParams()
		}
		return c.Redirect(Match.Show, match.ID)
	}

	return c.RenderText(fmt.Sprintf("Couldn't find a match with that id! (%d)", id))
}

// RemoveScore removes a match score from a match
func (c Match) RemoveScore(id int) revel.Result {
	ok, matchScore := models.FindMatchScore(id)
	if ok {
		matchScore.Delete()

		_, match := models.FindMatch(matchScore.MatchID)
		match.CalculateAll()
	}

	return c.Redirect(Match.Show, matchScore.MatchID)
}

// Delete deletes a match from the database
func (c Match) Delete(id int) revel.Result {
	ok, match := models.FindMatch(id)
	if ok {
		log.Print(match.Delete())
		return c.Redirect("/game/%d?group=%d", match.GameID, models.GroupIDFromGroupGameID(match.GroupGameID))
	}
	return c.RenderText(fmt.Sprintf("Couldn't find a match with that id! (%d)", id))
}
