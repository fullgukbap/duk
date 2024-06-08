package duk

import "net"

type (
	OnConnectHandler    = func(id string, conn net.Conn)
	OnDisconnectHandler = func(id string, err error)
)

type Hooks struct {
	app *App

	// Hooks
	onConnect    []OnConnectHandler
	onDisconnect []OnDisconnectHandler
}

func newHooks(app *App) *Hooks {
	return &Hooks{
		app:          app,
		onConnect:    make([]OnConnectHandler, 0),
		onDisconnect: make([]OnDisconnectHandler, 0),
	}
}

func (app *App) OnConnect(handler ...OnConnectHandler) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	app.hooks.onConnect = append(app.hooks.onConnect, handler...)
}

func (app *App) OnDisconnect(handler ...OnDisconnectHandler) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	app.hooks.onDisconnect = append(app.hooks.onDisconnect, handler...)
}
