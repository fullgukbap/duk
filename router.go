package duk

import "errors"

var (
	ErrNonEventExists = errors.New("event does not exist")
)

type Router struct {
	Handlers map[string]HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		Handlers: make(map[string]HandlerFunc),
	}
}

func (r *Router) Register(event string, handler HandlerFunc) {
	r.Handlers[event] = handler
}

func (r *Router) Path(event string) (HandlerFunc, error) {
	handler, exists := r.Handlers[event]
	if !exists {
		return nil, ErrNonEventExists
	}

	return handler, nil
}
