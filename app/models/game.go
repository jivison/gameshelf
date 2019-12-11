package models

import (
	"fmt"
	"log"
)

// Game is a board game
type Game struct {
	ID               int
	Title            string
	Year             int
	BggID            int
	ImgURL           string
	ComplexityRating float32
	Username         string `db:"user_name"`
	User             *User  `db:"-"`
}

func (g Game) String() string {
	return fmt.Sprintf("{ ID: %d | TITLE: %s | YEAR: %d | BGGID: %d | USERNAME: %s | IMGURL: %s | COMPLEXITYRATING: %f }", g.ID, g.Title, g.Year, g.BggID, g.Username, g.ImgURL, g.ComplexityRating)
}

// Matches returns all the matches associated with a game
func (g Game) Matches() []Match {
	var matches []Match
	dbmap.Select(&matches, "select * from matches where \"GameID\"=$1", g.ID)
	return matches
}

// Delete deletes a game from the database
func (g Game) Delete() bool {
	_, err := dbmap.Delete(&g)
	return (err != nil)
}

// Update updates the database with any changes to a Game
func (g Game) Update() error {
	_, err := dbmap.Update(&g)
	return err
}

// UnaddedGroups returns all the groups a game hasn't been added to yet
func (g Game) UnaddedGroups(username string) []Group {
	_, user := FindUser(username)

	usersGroups := user.Groups()

	var options []Group

	if len(usersGroups) == 0 {
		return options
	}

	var groupGames []GroupGame
	dbmap.Select(&groupGames, "select * from group_games where \"GameID\"=$1", g.ID)

	// Blacklist is the list of group (ID)s that the game has already been added to
	// Whitelist is the list of group (ID)s that the user is a part of
	var blacklist []int
	var whitelist []int

	for _, groupGame := range groupGames {
		blacklist = append(blacklist, groupGame.GroupID)
	}

	for _, group := range usersGroups {
		whitelist = append(whitelist, group.ID)
	}

	if len(blacklist) == 0 {
		dbmap.Select(&options, "select * from groups where \"ID\" in (:whitelist)", map[string]interface{}{
			"whitelist": whitelist,
		})
	} else {
		dbmap.Select(&options, "select * from groups where \"ID\" not in (:blacklist) and in (:whitelist)", map[string]interface{}{
			"blacklist": blacklist,
			"whitelist": whitelist,
		})
	}

	return options
}

// FindGame finds a game by its id
func FindGame(id int) (bool, *Game) {
	obj, err := dbmap.Get(Game{}, id)

	var game *Game

	if err != nil {
		log.Print("ERROR FindGame: ")
		log.Println(err)
	} else {
		game = obj.(*Game)
		_, game.User = FindUser(game.Username)
	}

	return (err == nil), game
}

// FindGameByTitle finds a game by its title
func FindGameByTitle(title, username string, storageVar *[]Game) {
	dbmap.Select(storageVar, "select * from games where \"Title\"=$1 and user_name=$1", title, username)
}

// CreateGame creates a game
func CreateGame(title string, year, bggID int, username, imgURL string, complexityRating float32) (bool, *Game) {
	game := &Game{
		Title:            title,
		Year:             year,
		BggID:            bggID,
		Username:         username,
		ImgURL:           imgURL,
		ComplexityRating: complexityRating,
	}

	err := dbmap.Insert(game)

	if err != nil {
		log.Print("ERROR: ")
		log.Println(err)
	}

	return (err == nil), game
}
