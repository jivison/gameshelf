package controllers

import (
	"gameshelf/app/models"

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
func (c User) Create(user models.User, verifyPassword string) revel.Result {
	user.HashedPassword, _ = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	models.InsertUser(user)
	c.Session["user"] = user.Username

	return c.Redirect(App.Index)
}
