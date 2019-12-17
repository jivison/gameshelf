package controllers

import (
	"gameshelf/app/models"
	"gameshelf/app/respond"

	"golang.org/x/crypto/bcrypt"

	"github.com/revel/revel"
)

// App is the default app controller
type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result  { return c.Render() }
func (c App) SignIn() revel.Result { return c.Render() }

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

type LoginParams struct {
	username string
	password string
}

// Login logs the user in
func (c App) Login(body LoginParams) revel.Result {
	user := c.getUser(body.username)
	errors := respond.NewErrors()

	if user != nil {
		err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(body.password))

		if err == nil {
			c.Session["user"] = body.username
			return respond.WithEntity(c, user)
		}
		errors.Add(err.Error())
	}
	errors.Add("Login failed")
	return respond.WithError(c, 400, *errors)
}

// SignOut signs the user out
func (c App) SignOut() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return respond.WithMessage(c, "Logged out")
}
