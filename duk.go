package duk

import (
	"encoding/json"
	"fmt"
	"io"
	"net"

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

type Duk struct {
	orchestration *Orchestration
	router        *Router
}

func New() *Duk {
	return &Duk{
		orchestration: NewOrchestration(),
		router:        NewRouter(),
	}
}

func (d *Duk) On(event string, handler HandlerFunc) {
	d.router.Register(event, handler)
}

func (d *Duk) Broadcast(data any) error {
	return d.orchestration.Broadcast(data)
}

func (d *Duk) Listen(port string) error {
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

		uuid := uuid.New().String()
		err = d.orchestration.Add(uuid, conn)
		if err != nil {
			return err
		}

		go d.requestHandler(uuid, conn)
	}
}

func (d *Duk) requestHandler(id string, conn net.Conn) {
	defer func() {
		conn.Close()
		d.orchestration.Remove(id)
	}()

	packet := Packet{}
	for {
		err := json.NewDecoder(conn).Decode(&packet)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		handler, err := d.router.Path(packet.Event)
		if err != nil {
			panic(err)
		}

		handler(NewCtx(packet.Payload))
	}
}
