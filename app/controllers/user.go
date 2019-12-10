package controllers

import (
	"fmt"
	"gameshelf/app/models"
	"regexp"

	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

// User Controller handles signing up
type User struct {
	*revel.Controller
}

// SignUp renders the sign up page
func (c User) SignUp() revel.Result {
	return c.Render()
}

// Create creates a persisted user
func (c User) Create(username, password, verifyPassword, firstName string) revel.Result {
	c.Validation.Required(username)
	c.Validation.MinSize(username, 4)
	c.Validation.MaxSize(username, 20)
	c.Validation.Match(username, regexp.MustCompile("^\\w*$"))

	c.Validation.Required(password)
	c.Validation.MinSize(password, 3)
	c.Validation.MaxSize(password, 30)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(User.SignUp)
	}

	user := models.User{
		Username:  username,
		Password:  password,
		FirstName: firstName,
	}

	user.HashedPassword, _ = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	models.InsertUser(user)

	c.Session["user"] = user.Username

	return c.Redirect(App.Index)
}

// Show displays a user and provides friending options
func (c User) Show(username string, from string) revel.Result {
	ok, user := models.FindUser(username)
	if ok == true {

		currentUser, _ := c.Session.Get("user")

		status := user.FriendStatus(currentUser.(string))

		return c.Render(user, status)
	}
	if from == "search" {
		c.Flash.Error("Couldn't find that user!")
		c.FlashParams()
		return c.Redirect(User.FindFriend)
	}
	return c.RenderText(fmt.Sprintf("Couldn't find a user with that username! (%s)", username))
}

// FindFriend displays the page to make a friend
func (c User) FindFriend() revel.Result {
	username, _ := c.Session.Get("user")
	ok, user := models.FindUser(username.(string))
	if !ok {
		return c.RenderText("You need to be signed in to access this!")
	}
	friends := user.Friends()
	pendingRequests := user.PendingFriendRequests()
	sentRequests := user.SentFriendRequests()

	return c.Render(friends, pendingRequests, sentRequests)
}

// AddFriend creates a friend request
func (c User) AddFriend(username string) revel.Result {
	sourceUsername, _ := c.Session.Get("user")
	models.CreateFriend(sourceUsername.(string), username)
	return c.Redirect(User.Show, username)
}

// AcceptRequest accepts a pending request
func (c User) AcceptRequest(username string) revel.Result {
	sourceUsername, _ := c.Session.Get("user")
	models.FindAndAcceptRequest(username, sourceUsername.(string))
	return c.Redirect(User.Show, username)
}
