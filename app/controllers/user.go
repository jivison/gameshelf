package controllers

import (
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
func (c User) Create(username, password, verifyPassword string) revel.Result {
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
		Username: username,
		Password: password,
	}

	user.HashedPassword, _ = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	models.InsertUser(user)
	c.Session["user"] = user.Username

	return c.Redirect(App.Index)
}
