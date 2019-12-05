package controllers

import (
	"fmt"
	"gameshelf/app/models"
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
func (c Match) Show(gameid, id int) revel.Result {
	ok, match := models.FindMatch(id)

	if ok {
		scores := match.MatchScores()
		return c.Render(match, scores)
	}

	return c.RenderText(fmt.Sprintf("Couldn't find a match with that id! (%d)", id))
}

func (c Match) AddScore(gameid, id int, playerUserName string, baseScore float32, isWinner bool) revel.Result {
	ok, game := models.FindGame(gameid)
	if ok {
		ok, match := models.FindMatch(id)
		if ok {
			models.CreateMatchScore(*match, *game, playerUserName, baseScore, isWinner)
		}

		return c.RenderText(fmt.Sprintf("Couldn't find a match with that id! (%d)", id))
	}

	return c.RenderText(fmt.Sprintf("Couldn't find a game with that id! (%d)", gameid))
}
