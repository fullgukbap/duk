package duk

import (
	"encoding/json"
	"errors"
	"net"
	"slices"
	"sync"
)

var (
	ErrIDExists    = errors.New("id that already exists")
	ErrNonIDExists = errors.New("id does not exists")
)

// Orchestration는 TCP Conn들을 관리하는 구조체 입니다.
type Orchestration struct {
	Conns map[string]net.Conn
	mu    sync.Mutex
}

func NewOrchestration() *Orchestration {
	return &Orchestration{
		Conns: make(map[string]net.Conn),
		mu:    sync.Mutex{},
	}
}

func (o *Orchestration) Add(id string, newConn net.Conn) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	_, exists := o.Conns[id]
	if exists {
		return ErrIDExists
	}

	o.Conns[id] = newConn
	return nil
}

func (o *Orchestration) Get(id string) (net.Conn, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	conn, exists := o.Conns[id]
	if !exists {
		return nil, ErrNonIDExists
	}

	return conn, nil
}

func (o *Orchestration) Remove(id string) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	_, exists := o.Conns[id]
	if !exists {
		return ErrNonIDExists
	}

	delete(o.Conns, id)
	return nil
}

func (o *Orchestration) Emit(id string, data any) error {
	conn, err := o.Get(id)
	if err != nil {
		return err
	}

	err = json.NewDecoder(conn).Decode(data)
	if err != nil {
		return err
	}

	return nil
}

func (o *Orchestration) Broadcast(data any, ignores ...string) error {
	for id, conn := range o.Conns {
		if slices.Contains(ignores, id) {
			continue
		}

		err := json.NewDecoder(conn).Decode(data)
		if err != nil {
			return err
		}
	}

	return nil
}
