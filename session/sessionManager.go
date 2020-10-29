package session

import (
	"sync"

	"github.com/gorilla/websocket"
)

var (
	managerOnce sync.Once
	managerInst *Manager
)

type Manager struct {
	sync.RWMutex
	sessions map[*websocket.Conn]*Session
}

func NewSessionManager() *Manager {
	return &Manager{
		sessions: make(map[*websocket.Conn]*Session),
	}
}

func ManagerInst() *Manager {
	managerOnce.Do(func() {
		managerInst = NewSessionManager()
	})
	return managerInst
}

func (object *Manager) Add(session *Session) *Session {
	object.Lock()
	defer object.Unlock()
	object.sessions[session.conn] = session
	return session
}

func (object *Manager) Remove(session *Session) *Session {
	object.Lock()
	defer object.Unlock()
	delete(object.sessions, session.conn)
	return session
}
