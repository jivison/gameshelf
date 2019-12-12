package controllers

import (
	"fmt"
	"gameshelf/app/models"
	"math"
	"sort"
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
func (c Game) Show(id int, group int) revel.Result {
	if ok, game := models.FindGame(id); ok {
		var matches []models.Match

		username, _ := c.Session.Get("user")
		groupChoices := game.UnaddedGroups(username.(string))

		if group != 0 {
			_, group := models.FindGroup(group)
			matches = game.Matches(group.ID)
			sort.SliceStable(matches, func(i, j int) bool {
				return matches[i].DatePlayed.Unix() > matches[j].DatePlayed.Unix()
			})
			return c.Render(game, matches, groupChoices, group)
		}

		return c.Render(game, groupChoices, matches)
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
func (c Game) Create(title, imgURL string, year, bggID int, complexityRating float64) revel.Result {

	username, _ := c.Session.Get("user")

	c.Log.Infof("{ TITLE: %s | YEAR: %d | BGGID: %d | USERNAME: %s | IMGURL: %s | COMPLEXITYRATING: %f}", title, year, bggID, username.(string), imgURL, complexityRating)

	if username != nil {
		c.Validation.Required(validateUniqueGame(title, username.(string), year)).Message("Title can't match another game with the same year")
	}

	c.Validation.Required(title).Message("A game must have a title")

	c.Validation.Required(username != nil).Message("You must be signed in to create a game")

	c.Validation.Range(year, 0, 2050).Message("Year must be realistic (0-2050)")

	c.Validation.Range(bggID, 0, math.MaxInt32).Message("Board Game Geek ID must be larger than 0")

	c.Validation.RangeFloat(complexityRating, 0, 5).Key("complexityRating").Message("Complexity Rating must be in between 0 and 5!")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Game.New)
	}

	if username != nil {
		if ok, game := models.CreateGame(title, year, bggID, username.(string), imgURL, float32(complexityRating)); ok {
			return c.Redirect(Game.Show, strconv.Itoa(game.ID))
		}
	}

	c.FlashParams()
	return c.Redirect(Game.New)
}

// Index lists every game
func (c Game) Index() revel.Result {
	username, _ := c.Session.Get("user")

	_, user := models.FindUser(username.(string))

	games := user.Games()
	groupChoices := user.Groups()

	sort.SliceStable(games, func(i, j int) bool {
		return games[i].Title < games[j].Title
	})

	return c.Render(games, groupChoices)
}

// Update updates a game in the database
func (c Game) Update(id int, title, imgURL string, year, bggID int, complexityRating float32) revel.Result {
	_, game := models.FindGame(id)

	if game.ComplexityRating != complexityRating {
		game.ComplexityRating = complexityRating
		game.Update()
		game.UpdateAllMatches()
	}

	game.Title = title
	game.Year = year
	game.BggID = bggID
	game.ImgURL = imgURL
	game.Update()
	return c.Redirect(fmt.Sprintf("/game/%d", game.ID))
}

// Delete deletes a game in the database
func (c Game) Delete(id int) revel.Result {
	_, game := models.FindGame(id)
	game.Delete()
	return c.Redirect(Game.Index)
}

// Edit renders the edit game page
func (c Game) Edit(id int) revel.Result {
	_, game := models.FindGame(id)
	return c.Render(game)
}

// AddAllToGroup adds every single one of a user's games to a group
func (c Game) AddAllToGroup(groupID int) revel.Result {
	ok, group := models.FindGroup(groupID)
	if ok {
		username, _ := c.Session.Get("user")
		group.AddAllGames(username.(string))
	}
	return c.Redirect(Game.Index)
}

// AddToGroup adds a game to a group, creating a GroupGame in the process
func (c Game) AddToGroup(id, groupID int) revel.Result {
	ok, _ := models.FindGroup(groupID)
	if ok {
		models.CreateGroupGame(groupID, id)
		return c.Redirect(Game.Show, id)
	}
	return c.RenderText("Couldn't find a group with that id! (%d)", groupID)
}
