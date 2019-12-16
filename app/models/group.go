package models

import (
	"log"
	"sort"
)

// Group is a group of users
type Group struct {
	ID   int
	Name string
}

// Scoreboard type represents a group's scoreboard
type Scoreboard struct {
	Scores []ScoreboardScore
}

// Sort sorts the scoreboard
func (s Scoreboard) Sort() {
	sort.SliceStable(s.Scores, func(i, j int) bool {
		return s.Scores[i].AvgScore > s.Scores[j].AvgScore
	})
}

// ScoreboardScore holds the individual player scores in a scoreboard
type ScoreboardScore struct {
	AggScore    float32
	AvgScore    float32
	Count       int
	DisplayName string
}

// Members returns a list of users in a group
func (g Group) Members() []User {
	var groupMembers []GroupMember
	dbmap.Select(&groupMembers, "select * from group_members where \"GroupID\"=$1 and \"Pending\"='f'", g.ID)

	var usernames []string

	for _, groupMember := range groupMembers {
		usernames = append(usernames, groupMember.Username)
	}

	var users []User

	dbmap.Select(&users, "select * from users where \"Username\" in (:usernames)", map[string]interface{}{
		"usernames": usernames,
	})

	return users
}

// SentInvitations returns all sent invitations from the group
func (g Group) SentInvitations() []GroupMember {
	var invitations []GroupMember
	dbmap.Select(&invitations, "select * from group_members where \"GroupID\"=$1 and \"Pending\"='t'", g.ID)
	return invitations
}

// Games returns all the games associated with a group
func (g Group) Games() []Game {
	var groupGames []GroupGame
	dbmap.Select(&groupGames, "select * from group_games where \"GroupID\"=$1", g.ID)

	var gameIDs []int

	for _, groupGame := range groupGames {
		gameIDs = append(gameIDs, groupGame.GameID)
	}

	var games []Game

	dbmap.Select(&games, "select * from games where \"ID\" in (:gameids)", map[string]interface{}{
		"gameids": gameIDs,
	})
	return games
}

// AddAllGames adds all of a user's games a group
func (g Group) AddAllGames(username string) {
	var games []Game

	var blacklist []int

	for _, game := range g.Games() {
		blacklist = append(blacklist, game.ID)
	}

	dbmap.Select(&games, "select * from games where \"ID\" not in (:blacklist) and user_name=:username ", map[string]interface{}{
		"username":  username,
		"blacklist": blacklist,
	})

	for _, game := range games {
		CreateGroupGame(g.ID, game.ID)
	}
}

// GroupGames returns all the groupGames associated with a group
func (g Group) GroupGames() []GroupGame {
	var groupGames []GroupGame
	dbmap.Select(&groupGames, "select * from group_games where \"GroupID\"=$1", g.ID)
	return groupGames
}

// MatchScores returns all the MatchScores associated with a group
func (g Group) MatchScores() []MatchScore {
	groupGames := g.GroupGames()

	var whitelist []int

	for _, gg := range groupGames {
		whitelist = append(whitelist, gg.ID)
	}

	var matches []Match
	dbmap.Select(&matches, "select * from matches where \"GroupGameID\" in (:whitelist)", map[string]interface{}{
		"whitelist": whitelist,
	})

	var matchScores []MatchScore

	for _, m := range matches {
		matchScores = append(matchScores, m.MatchScores()...)
	}

	return matchScores
}

// Scoreboard returns a scoreboard of players' scores
func (g Group) Scoreboard() Scoreboard {
	tempScoreboard := make(map[string]*ScoreboardScore)

	for _, ms := range g.MatchScores() {
		if _, ok := tempScoreboard[ms.PlayerUserName]; !ok {
			tempScoreboard[ms.PlayerUserName] = &ScoreboardScore{
				DisplayName: ms.PlayerDisplayName,
			}
		}
		tempScoreboard[ms.PlayerUserName].AggScore = ms.FinalScore
		tempScoreboard[ms.PlayerUserName].Count++
	}

	scoreboard := Scoreboard{}

	for player, scores := range tempScoreboard {
		tempScoreboard[player].AvgScore = scores.AggScore / float32(scores.Count)
		scoreboard.Scores = append(scoreboard.Scores, *tempScoreboard[player])
	}

	return scoreboard
}

// CreateGroup creates a group in the database
func CreateGroup(name string) (bool, *Group) {
	group := &Group{
		Name: name,
	}

	err := dbmap.Insert(group)

	if err != nil {
		log.Print("ERROR CreateGroup:")
		log.Print(err)
	}

	return (err == nil), group
}

// FindGroup finds a group by its ID
func FindGroup(id int) (bool, *Group) {
	var groups []Group

	dbmap.Select(&groups, findQstring("groups"), id)

	if len(groups) > 0 {
		return true, &groups[0]
	}
	return false, &Group{}
}
