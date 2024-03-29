package controllers

import (
	"gameshelf/app/models"

	"golang.org/x/crypto/bcrypt"

	"github.com/revel/revel"
)

// App is the default app controller
type App struct {
	*revel.Controller
}

// Index renders the index page
func (c App) Index() revel.Result {
	username, _ := c.Session.Get("user")
	fulluser, _ := c.Session.Get("fulluser")
	return c.Render(username, fulluser)
}

func (c App) connected() *models.User {
	if c.ViewArgs["user"] != nil {
		return c.ViewArgs["user"].(*models.User)
	}
	if username, ok := c.Session["user"]; ok {
		return c.getUser(username.(string))
	}

	return nil
}

func (c App) getUser(username string) (user *models.User) {
	user = &models.User{}
	c.Session.GetInto("fulluser", user, false)

	if user.Username == username {
		return user
	}

	ok, user := models.FindUser(username)

	c.Log.Info(user.String())

	if !ok {
		c.Log.Errorf("Couldn't find user with the username: %s", username)
		return nil
	}

	c.Session["fulluser"] = user
	return user

}

// Login logs the user in (does not display the page)
func (c App) Login(username, password string) revel.Result {
	user := c.getUser(username)
	if user != nil {
		err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))

		if err == nil {
			c.Session["user"] = username
			c.Flash.Success("Welcome, " + username)
			return c.Redirect(App.Index)

		}
	}
	c.Flash.Out["username"] = username
	c.Flash.Error("Login Failed")
	return c.Redirect(App.SignIn)
}

// SignIn displays the signin page
func (c App) SignIn() revel.Result {
	return c.Render()
}

// SignOut signs the user out
func (c App) SignOut() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(App.Index)
}
