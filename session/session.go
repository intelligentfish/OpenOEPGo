package session

import (
	"fmt"
	"os"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"openOEP/pb"
	"openOEP/user"
)

// Session
type Session struct {
	user *user.User
	conn *websocket.Conn
}

func NewSession(conn *websocket.Conn) *Session {
	return &Session{conn: conn}
}

func (object *Session) Conn() *websocket.Conn {
	return object.conn
}

func (object *Session) User() *user.User {
	return object.user
}

func (object *Session) Loop() {
	var mt int
	var raw []byte
	var err error
	defer object.conn.Close()
	for {
		mt, raw, err = object.conn.ReadMessage()
		if nil != err {
			fmt.Fprintln(os.Stderr, err)
			break
		}
		switch mt {
		case websocket.TextMessage:
		case websocket.BinaryMessage:
			var req pb.Request
			if err = proto.Unmarshal(raw, &req); nil != err {
				fmt.Fprintln(os.Stderr, err)
				break
			}
			switch req.RequestType {
			case pb.Request_NALU:
				var nal pb.NALRequest
				if err = proto.Unmarshal(req.Body, &nal); nil != err {
					fmt.Fprintln(os.Stderr, err)
					break
				}
				switch nal.NalType {
				case pb.NALRequest_VPS:
				case pb.NALRequest_SPS:
				case pb.NALRequest_PPS:
				}
			}
		}
	}
}
