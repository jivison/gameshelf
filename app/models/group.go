package models

import (
	"log"
)

// Group is a group of users
type Group struct {
	ID   int
	Name string
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
	obj, err := dbmap.Get(Group{}, id)
	group := obj.(*Group)

	return (err == nil), group
}
