package controllers

import (
	"github.com/revel/revel"
)

// BaseController is the base type for all controllers (and holds generic actions)
type BaseController struct {
	*revel.Controller
}
