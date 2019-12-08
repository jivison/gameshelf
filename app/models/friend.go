package models

import "log"

// Friend represents the join table connecting users to itself
type Friend struct {
	ID               int
	FrienderUsername string
	FriendedUsername string
	Pending          bool
}

// Update updates a friend in the database
func (f Friend) Update() {
	dbmap.Update(&f)
}

// AcceptRequest de-pends a Friend in the database
func (f Friend) AcceptRequest() {
	f.Pending = true
	f.Update()
}

// CreateFriend creates a friend in the database
func CreateFriend(sourceUsername, targetUsername string) (bool, *Friend) {
	friend := &Friend{
		FrienderUsername: sourceUsername,
		FriendedUsername: targetUsername,
		Pending:          true,
	}

	err := dbmap.Insert(friend)

	if err != nil {
		log.Print("ERROR: ")
		log.Println(err)
	}

	return (err == nil), friend
}
