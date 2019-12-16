package respond

import "github.com/revel/revel"

type controller interface {
	RenderJSON(o interface{}) revel.Result
}

// WithError responds with an error
func WithError(c controller, status int, errors JSONErrors) revel.Result {
	return c.RenderJSON(ErrorResponse{
		JSONResponse{Status: status},
		errors,
	})
}

// WithMessage responds with an error
func WithMessage(c controller, messages ...string) revel.Result {
	return c.RenderJSON(MessageResponse{
		JSONResponse{Status: 200},
		messages,
	})
}

// WithEntity responds with an error
func WithEntity(c controller, entity interface{}) revel.Result {
	return c.RenderJSON(SingleEntityResponse{
		JSONResponse{Status: 200},
		entity,
	})
}

// WithEntities responds with an error
func WithEntities(c controller, entities ...interface{}) revel.Result {
	return c.RenderJSON(MultipleEntityResponse{
		JSONResponse{Status: 200},
		entities,
	})
}
