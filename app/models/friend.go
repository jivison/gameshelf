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
func (f Friend) AcceptRequest() bool {
	f.Pending = false
	f.Update()
	newFriend := &Friend{
		FrienderUsername: f.FriendedUsername,
		FriendedUsername: f.FrienderUsername,
		Pending:          false,
	}

	err := dbmap.Insert(newFriend)
	return (err == nil)
}

// CreateFriend creates a friend in the database
func CreateFriend(sourceUsername, targetUsername string) (bool, *Friend) {

	if !validateUniqueFriend(sourceUsername, targetUsername) {
		return false, nil
	}

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

func validateUniqueFriend(source, target string) bool {
	var friends []Friend
	dbmap.Select(&friends, "select * from friends where (\"FrienderUsername\"=:source and \"FriendedUsername\"=:target) OR (\"FriendedUsername\"=:source and \"FrienderUsername\"=:target)", map[string]interface{}{
		"source": source,
		"target": target,
	})
	return (len(friends) == 0)
}

// FindAndAcceptRequest finds and accepts a pending friend request
func FindAndAcceptRequest(friender, friended string) bool {
	log.Print(friender, " - ", friended)

	var friends []Friend
	dbmap.Select(&friends, "select * from friends where \"FrienderUsername\"=:friender and \"FriendedUsername\"=:friended", map[string]interface{}{
		"friender": friender,
		"friended": friended,
	})

	if len(friends) != 0 {
		request := friends[0]
		if request.Pending {
			return request.AcceptRequest()
		}
		return false
	}
	return false
}
