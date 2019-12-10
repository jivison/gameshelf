package models

import (
	"fmt"
	"log"
	"strings"
)

// User is a user
type User struct {
	Username       string `db:",primarykey"`
	Password       string `db:"-"`
	FirstName      string
	HashedPassword []byte
}

// FriendStatus is the relationship of two users
type FriendStatus struct {
	// Status can either be 'friends', 'sentRequest', 'receivedRequest', 'notFriends'
	Status  string
	Pending bool
	// Blocked bool // Coming some time in the future
}

func (u User) String() string {
	return fmt.Sprintf("{ USERNAME: %s | HASHEDPASSWORD: %s }", u.Username, u.HashedPassword)
}

// FriendStatus returns the relationship between users
func (u User) FriendStatus(username string) FriendStatus {
	var friends []Friend
	dbmap.Select(&friends, "select * from friends where \"FrienderUsername\"=:friender and \"FriendedUsername\"=:friended", map[string]interface{}{
		"friended": u.Username,
		"friender": username,
	})

	log.Print(friends, u.Username, username)

	if len(friends) != 0 {
		if friends[0].Pending {
			return FriendStatus{
				Status:  "sentRequest",
				Pending: true,
			}
		}
		return FriendStatus{
			Status:  "friends",
			Pending: false,
		}
	}

	var receivedRequest []Friend
	dbmap.Select(&receivedRequest, "select * from friends where \"FriendedUsername\"=:friended and \"FrienderUsername\"=:friender", map[string]interface{}{
		"friended": username,
		"friender": u.Username,
	})

	if len(receivedRequest) != 0 {
		return FriendStatus{
			Status:  "receivedRequest",
			Pending: true,
		}
	}

	return FriendStatus{
		Status:  "notFriends",
		Pending: false,
	}
}

// AllFriends returns a list of all the friends associated with that user (pending or not)
func (u User) AllFriends() []Friend {
	var friends []Friend
	dbmap.Select(&friends, "select * from friends where \"FrienderUsername\"=$1", u.Username)
	return friends
}

// Friends returns a list of all the friends associated with that user (that aren't pending)
func (u User) Friends() []Friend {
	var friends []Friend
	dbmap.Select(&friends, "select * from friends where \"FrienderUsername\"=$1 AND \"Pending\"='f'", u.Username)
	return friends
}

// SentFriendRequests returns a list of all the friends associated with that user (that are pending)
// (ie. the ones the user has sent that haven't been accepted)
func (u User) SentFriendRequests() []Friend {
	var friends []Friend
	dbmap.Select(&friends, "select * from friends where \"FrienderUsername\"=$1 AND \"Pending\"='t'", u.Username)
	return friends
}

// PendingFriendRequests returns a list of all pending friends that need to be accepted
func (u User) PendingFriendRequests() []Friend {
	var friends []Friend
	dbmap.Select(&friends, "select * from friends where \"FriendedUsername\"=$1 AND \"Pending\"='t'", u.Username)
	return friends
}

// Games returns all the games associated with a user
func (u User) Games() []Game {
	var games []Game
	dbmap.Select(&games, "select * from games where user_name=$1", u.Username)
	return games
}

// FindUser finds a user by its username
func FindUser(username string) (bool, *User) {
	var users []User
	_, err := dbmap.Select(&users, "select * from users where \"Username\"=$1", strings.TrimSpace(username))

	if err != nil {
		log.Print("ERROR FindUser: ")
		log.Println(err)
	}

	if len(users) < 1 {
		log.Printf("INFO: FindUser couldn't find a user with username %s", username)
		return false, &User{}
	}

	user := users[0]

	return (err == nil), &user
}

// InsertUser inserts an already 'built' user
func InsertUser(user User) bool {
	err := dbmap.Insert(&user)

	if err != nil {
		log.Print("ERROR: ")
		log.Println(err)
	}

	return err == nil

}

// FriendableUsers returns a list of all the users who aren't already friended, and who haven't
// sent a friend request
func FriendableUsers(username string) []User {
	_, user := FindUser(username)
	friends := user.AllFriends()

	var users []User
	var blacklist []string

	blacklist = append(blacklist, username)

	for _, friend := range friends {
		blacklist = append(blacklist, friend.FriendedUsername)
	}

	for _, friend := range user.PendingFriendRequests() {
		blacklist = append(blacklist, friend.FrienderUsername)
	}

	dbmap.ExpandSliceArgs = true
	dbmap.Select(&users, "SELECT \"Username\" FROM users WHERE \"Username\" NOT IN (:Blacklist)", map[string]interface{}{
		"Blacklist": blacklist,
	})

	return users
}
