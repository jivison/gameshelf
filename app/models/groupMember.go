package models

// GroupMember joins the user and group tables in a many to many relationship
type GroupMember struct {
	ID       int
	GroupID  int
	Username string
	Pending  bool
}

// Update updates a group member in the database
func (gm GroupMember) Update() {
	dbmap.Update(&gm)
}

// AcceptInvitation de-pends a group member
func (gm GroupMember) AcceptInvitation() {
	gm.Pending = false
	gm.Update()
}

// CreateGroupMember creates a group member in the database
func CreateGroupMember(groupID int, username string) (bool, *GroupMember) {
	groupMember := &GroupMember{
		GroupID:  groupID,
		Username: username,
		Pending:  false,
	}

	err := dbmap.Insert(groupMember)

	return (err == nil), groupMember
}

// CreateGroupInvitation creates a pending group member in the database
func CreateGroupInvitation(groupID int, username string) (bool, *GroupMember) {
	groupMember := &GroupMember{
		GroupID:  groupID,
		Username: username,
		Pending:  true,
	}

	err := dbmap.Insert(groupMember)

	return (err == nil), groupMember
}

// FindGroupInvitation finds a group invitation from a groupID and a username
func FindGroupInvitation(groupID int, username string) *GroupMember {
	var invitations []GroupMember
	dbmap.Select(&invitations, "select * from group_members where \"Username\"=:username and \"GroupID\"=:groupid AND \"Pending\"='t'", map[string]interface{}{
		"username": username,
		"groupid":  groupID,
	})

	if len(invitations) > 0 {
		return &invitations[0]
	}
	return &GroupMember{}
}
