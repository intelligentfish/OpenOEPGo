package room

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"openOEP/session"
)

// Room
type Room struct {
	sync.RWMutex
	Id              int                                           `json:"id"`
	users           map[int]*session.Session                      `json:"users"`
	sessions        map[*websocket.Conn]*session.Session          `json:"sessions"`
	offlineCallback func(roomId int, sessions []*session.Session) `json:"-"`
}

// NewRoom
func NewRoom(id int,
	sessions []*session.Session,
	offlineCallback func(roomId int, users []*session.Session)) *Room {
	object := &Room{
		Id:              id,
		sessions:        make(map[*websocket.Conn]*session.Session),
		offlineCallback: offlineCallback,
	}
	for _, session := range sessions {
		object.sessions[session.Conn()] = session
	}
	return object
}

// Broadcast
func (object *Room) Broadcast(messageType int,
	raw []byte,
	timeout time.Duration) *Room {
	object.RLock()
	defer object.RUnlock()
	var err error
	var sessions []*session.Session
	for _, session := range object.sessions {
		session.Conn().SetWriteDeadline(time.Now().Add(timeout))
		if err = session.Conn().WriteMessage(messageType, raw);
			nil != err {
			sessions = append(sessions, session)
			fmt.Fprintln(os.Stderr, err)
		}
	}
	if nil != sessions && 0 < len(sessions) {
		object.offlineCallback(object.Id, sessions)
	}
	return object
}

// Send
func (object *Room) Send(messageType int,
	raw []byte,
	userId int,
	timeout time.Duration) (found bool, err error) {
	object.RLock()
	defer object.RUnlock()
	var session *session.Session
	if session, found = object.users[userId]; found {
		session.Conn().SetWriteDeadline(time.Now().Add(timeout))
		err = session.Conn().WriteMessage(messageType, raw)
	}
	return
}

// Join
func (object *Room) Join(session *session.Session) {
	object.Lock()
	defer object.Unlock()
	object.users[session.User().Id] = session
}

// Leave
func (object *Room) Leave(messageType int,
	raw []byte,
	userId int,
	timeout time.Duration) (found bool, err error) {
	found, err = object.Send(messageType, raw, userId, timeout)
	object.Lock()
	defer object.Unlock()
	s := object.users[userId]
	delete(object.users, userId)
	object.offlineCallback(object.Id, []*session.Session{s})
	return
}
