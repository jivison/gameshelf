package controllers

import (
	"fmt"
	"gameshelf/app/models"

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

// Create creates a new game in the db
func (c Game) Create(title string, year, bggID int) revel.Result {
	if ok, game := models.CreateGame(title, year, bggID); ok {
		return c.Redirect(fmt.Sprintf("/game/%d", game.ID))
	}

	c.FlashParams()
	return c.Redirect(Game.New)
}
