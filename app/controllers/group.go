package controllers

import (
	"fmt"
	"gameshelf/app/models"

	"github.com/revel/revel"
)

// Group controller handles the creation, deletion, and other actions related to groups
type Group struct {
	*revel.Controller
}

// New displays the page for a new group
func (c Group) New() revel.Result {
	return c.Render()
}

// Create creates a group in the database
func (c Group) Create(name string) revel.Result {
	c.Validation.Required(name).Message("A name is required")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Group.New)
	}

	if ok, group := models.CreateGroup(name); ok {
		username, _ := c.Session.Get("user")
		models.CreateGroupMember(group.ID, username.(string))

		return c.Redirect(Group.Show, group.ID)
	}

	c.FlashParams()
	return c.Redirect(Group.New)
}

// Show displays a group
func (c Group) Show(id int) revel.Result {
	ok, group := models.FindGroup(id)
	username, _ := c.Session.Get("user")
	username = username.(string)
	if ok {
		members := group.Members()
		sentInvitations := group.SentInvitations()
		games := group.Games()
		return c.Render(group, members, sentInvitations, games, username)
	}
	return c.RenderText(fmt.Sprintf("Couldn't find a group with that ID! (%d)", id))
}

// Index displays a list of groups
func (c Group) Index() revel.Result {
	username, _ := c.Session.Get("user")

	_, user := models.FindUser(username.(string))

	groups := user.Groups()
	invitations := user.PendingGroupInvitations()

	return c.Render(groups, invitations)
}

// SendInvitation sends a group invitaation to a user
func (c Group) SendInvitation(id int, username string) revel.Result {
	ok, _ := models.FindGroup(id)
	if ok {
		models.CreateGroupInvitation(id, username)
		return c.Redirect(Group.Show, id)
	}
	return c.RenderText(fmt.Sprintf("Couldn't find a group with that id! (%d)", id))
}

// UnsendInvitation deletes a pending group invitation
func (c Group) UnsendInvitation(id int, username string) revel.Result {
	ok, _ := models.FindGroup(id)
	if ok {
		models.DestroyGroupInvitation(id, username)
		return c.Redirect(Group.Show, id)
	}
	return c.RenderText(fmt.Sprintf("Couldn't find a group with that id! (%d)", id))
}

// AcceptInvitation accepts an invitation to a group
func (c Group) AcceptInvitation(groupID int) revel.Result {
	username, _ := c.Session.Get("user")
	ok, _ := models.FindGroup(groupID)
	if ok {
		invite := models.FindGroupInvitation(groupID, username.(string))
		invite.AcceptInvitation()
		return c.Redirect(Group.Show, groupID)
	}
	return c.RenderText(fmt.Sprintf("Couldn't find a group with that id! (%d)", groupID))
}
