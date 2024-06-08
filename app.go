package duk

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/google/uuid"
)

const (
	Version = "1.0.0-beta"
	banner  = "" +
		" \x1b[1;32m ___       _    \n" +
		" \x1b[1;32m| . \\ _ _ | |__\n" +
		" \x1b[1;32m| | || | || / / \n" +
		" \x1b[1;32m|___/\\___||_\\_\\ \n" +
		" \x1b[1;30m%s\x1b[1;32m%v\x1b[0000m\n\n"
)

type App struct {
	mutex sync.Mutex

	router        *Router
	orchestration *Orchestration
	hooks         *Hooks
}

func New() *App {
	app := &App{
		mutex: sync.Mutex{},
	}

	app.router = newRouter(app)
	app.orchestration = newOrchestration(app)
	app.hooks = newHooks(app)

	return app
}

func (app *App) Listen(port string) error {
	fmt.Printf(banner, Version, fmt.Sprintf(" port%s", port))

	ln, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		id := uuid.New().String()
		err = app.orchestration.Add(id, conn)
		if err != nil {
			return err
		}

		for _, handler := range app.hooks.onConnect {
			handler(id, conn)
		}

		go app.requestHandler(id, conn)
	}
}

func (app *App) requestHandler(id string, conn net.Conn) {
	defer func() {
		conn.Close()
		app.orchestration.Remove(id)
	}()

	packet := Packet{}
	for {
		err := json.NewDecoder(conn).Decode(&packet)
		if err != nil {
			if err == io.EOF {
				for _, handler := range app.hooks.onDisconnect {
					handler(id, err)
				}
				break
			}
			panic(err)
		}

		handler, err := app.router.match(packet.Event)
		if err != nil {
			panic(err)
		}

		handler(NewCtx(packet.Payload))
	}
}
