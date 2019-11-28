package controllers

import (
	"gameshelf/app/models"
	"math"
	"strconv"

	"github.com/revel/revel"
)

// Game controller handles all actions regarding games
type Game struct {
	*revel.Controller
}

// New serves a form to create a new game
func (c Game) New() revel.Result {
	return c.Render()
}

// Show displays the details of a single game
func (c Game) Show(id string) revel.Result {
	numID, err := strconv.Atoi(id)
	if err != nil {
		return c.RenderError(err)
	}
	if ok, game := models.FindGame(numID); ok {
		return c.Render(game)
	}
	return c.RenderText("Couldn't find a game with that id!")
}

func validateUniqueGame(title, username string, year int) bool {
	var games []models.Game

	models.FindGameByTitle(title, username, &games)

	for _, persistedGame := range games {
		if persistedGame.Year != 0 && persistedGame.Year == year {
			return false
		}
	}

	return true

}

// Create creates a new game in the db
func (c Game) Create(title string, year, bggID int) revel.Result {

	username, _ := c.Session.Get("user")

	c.Log.Infof("TITLE: %s | YEAR: %d | BGGID: %d", title, year, bggID)

	if username != nil {
		c.Validation.Required(validateUniqueGame(title, username.(string), year)).Message("Title can't match another game with the same year")
	}

	c.Validation.Required(title).Message("A game must have a title")

	c.Validation.Required(username != nil).Message("You must be signed in to create a game")

	c.Validation.Range(year, 0, 2050).Message("Year must be realistic (0-2050)")

	c.Validation.Range(bggID, 0, math.MaxInt32).Message("Board Game Geek ID must be larger than 0")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Game.New)
	}

	if username, err := c.Session.Get("user"); err == nil {
		if ok, game := models.CreateGame(title, year, bggID, username.(string)); ok {
			return c.Redirect(Game.Show, strconv.Itoa(game.ID))
		}
	}

	c.FlashParams()
	return c.Redirect(Game.New)
}
