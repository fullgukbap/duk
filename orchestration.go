package duk

import (
	"encoding/json"
	"errors"
	"net"
	"slices"
)

var (
	ErrIDExists    = errors.New("id that already exists")
	ErrNonIDExists = errors.New("id does not exists")
)

// Orchestration는 TCP Conn들을 관리하는 구조체 입니다.
type Orchestration struct {
	app *App

	Conns map[string]net.Conn
}

func newOrchestration(app *App) *Orchestration {
	return &Orchestration{
		app:   app,
		Conns: make(map[string]net.Conn),
	}
}

func (o *Orchestration) add(id string, newConn net.Conn) error {
	o.app.mutex.Lock()
	defer o.app.mutex.Unlock()

	_, exists := o.Conns[id]
	if exists {
		return ErrIDExists
	}

	o.Conns[id] = newConn
	return nil
}

func (o *Orchestration) get(id string) (net.Conn, error) {
	o.app.mutex.Lock()
	defer o.app.mutex.Unlock()

	conn, exists := o.Conns[id]
	if !exists {
		return nil, ErrNonIDExists
	}

	return conn, nil
}

func (o *Orchestration) remove(id string) error {
	o.app.mutex.Lock()
	defer o.app.mutex.Unlock()

	_, exists := o.Conns[id]
	if !exists {
		return ErrNonIDExists
	}

	delete(o.Conns, id)
	return nil
}

func (o *Orchestration) broadcast(data any, ignores ...string) error {
	for id, conn := range o.Conns {
		if slices.Contains(ignores, id) {
			continue
		}

		err := json.NewEncoder(conn).Encode(data)
		if err != nil {
			return err
		}
	}

	return nil
}
