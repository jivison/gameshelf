package controllers

import (
	"gameshelf/app/models"
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

// Create creates a new game in the db
func (c Game) Create(title string, year, bggID int) revel.Result {
	if ok, game := models.CreateGame(title, year, bggID); ok {
		return c.Redirect(Game.Show, strconv.Itoa(game.ID))
	}

	c.FlashParams()
	return c.Redirect(Game.New)
}
