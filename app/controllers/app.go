package controllers

import (
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
