package controllers

import (
	"fmt"
	"gameshelf/app/models"
	"gameshelf/app/respond"
	"regexp"

	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

// User Controller handles signing up
type User struct {
	*revel.Controller
}

func (c User) SignUp() revel.Result { return c.Render() }

// Create creates a persisted user
func (c User) Create(username, password, verifyPassword, firstName string) revel.Result {

	c.Validation.Required(username).Message("Username is required")
	c.Validation.MinSize(username, 4).Message("Username must be at least 4 characters long")
	c.Validation.MaxSize(username, 40).Message("Username must be less than 40 characters")
	c.Validation.Match(username, regexp.MustCompile("^\\w*$")).Message("Username must not include special characters other than _ or .")

	c.Validation.Required(password).Message("Password is required")
	c.Validation.MinSize(password, 3).Message("Password must be longer than 2 characters")
	c.Validation.MaxSize(password, 40).Message("Password must be less than 40 characters")

	errors := respond.NewErrors()

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		errors.AddFromValidation(c.Validation.Errors)
		return respond.WithError(c, 422, *errors)
	}

	if firstName == "" {
		firstName = username
	}

	user := models.User{
		Username:  username,
		Password:  password,
		FirstName: firstName,
	}

	user.HashedPassword, _ = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	models.InsertUser(user)

	c.Session["user"] = user.Username

	return respond.WithEntity(c, user)
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
