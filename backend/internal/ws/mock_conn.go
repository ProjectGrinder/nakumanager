package ws

type MockConn struct {
	Messages []interface{}
}

func (m *MockConn) WriteJSON(v interface{}) error {
	m.Messages = append(m.Messages, v)
	return nil
}

func (m *MockConn) ReadMessage() (int, []byte, error) {
	return 1, []byte("test"), nil
}

func (m *MockConn) Close() error {
	return nil
}
