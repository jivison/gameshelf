package controllers

import (
	"gameshelf/app/models"

	"github.com/revel/revel"
)

// App is the default app controller
type App struct {
	*revel.Controller
}

// Index renders the index page
func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) AddUser() revel.Result {
	if user := c.connected(); user != nil {
		c.ViewArgs["user"] = user
	}
	return nil
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

	if !ok {
		c.Log.Errorf("Couldn't find user with the username: %s", username)
		return nil
	}

	c.Session["fulluser"] = user
	return

}
