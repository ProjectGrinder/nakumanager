package ws

import wsfiber "github.com/gofiber/contrib/websocket"

type WSConn interface {
	WriteJSON(v interface{}) error
	ReadMessage() (int, []byte, error)
	Close() error
}

type FiberWSConn struct {
	Conn *wsfiber.Conn
}

func (f *FiberWSConn) WriteJSON(v interface{}) error {
	return f.Conn.WriteJSON(v)
}

func (f *FiberWSConn) ReadMessage() (int, []byte, error) {
	return f.Conn.ReadMessage()
}

func (f *FiberWSConn) Close() error {
	return f.Conn.Close()
}
