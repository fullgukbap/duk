package duk

import "errors"

var (
	ErrNonEventExists = errors.New("event does not exist")
)

type Router struct {
	app *App

	handlers map[string]HandlerFunc
}

func newRouter(app *App) *Router {
	return &Router{
		app:      app,
		handlers: make(map[string]HandlerFunc),
	}
}

func (a *App) On(event string, handler HandlerFunc) {
	a.router.handlers[event] = handler
}

func (r *Router) match(event string) (HandlerFunc, error) {
	handler, exists := r.handlers[event]
	if !exists {
		return nil, ErrNonEventExists
	}

	return handler, nil
}
