package mediaServer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"openOEP/room"
	"openOEP/session"
)

// Server
type Server struct {
	port int `json:"port" toml:"port"`
}

// New
func New() *Server {
	return &Server{}
}

// onUserOffline
func (object *Server) onUserOffline(roomId int, sessions []*session.Session) {

}

// String
func (object *Server) String() string {
	raw, _ := json.Marshal(object)
	return string(raw)
}

// SetPort
func (object *Server) SetPort(port int) *Server {
	object.port = port
	return object
}

// Start
func (object *Server) Start(wg *sync.WaitGroup) (err error) {
	// add word room
	room.ManagerInst().AddRoom(room.NewRoom(0, nil, object.onUserOffline))
	wg.Add(1)
	go func() {
		defer wg.Done()
		var upgrader websocket.Upgrader
		http.HandleFunc("/stream", func(w http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(w, r, nil)
			if nil != err {
				fmt.Fprintln(os.Stderr, err)
				return
			}
			session.ManagerInst().Add(session.NewSession(conn)).Loop()
		})
		http.ListenAndServe(fmt.Sprintf(":%d", object.port), nil)
	}()
	return
}

// Stop
func (object *Server) Stop() (err error) {
	return
}
