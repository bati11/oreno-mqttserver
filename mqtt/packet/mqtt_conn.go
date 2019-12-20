package packet

import (
	"bufio"
	"net"

	"github.com/gorilla/websocket"
)

func NewMQTTConn(conn net.Conn) *MQTTConn {
	return &MQTTConn{conn: conn}
}

func NewMQTTConnWithWebSocket(websocketConn *websocket.Conn) *MQTTConn {
	return &MQTTConn{websocketConn: websocketConn}
}

type MQTTConn struct {
	websocketConn *websocket.Conn
	conn          net.Conn
}

func (c *MQTTConn) Close() error {
	if c.websocketConn != nil {
		return c.websocketConn.Close()
	}
	return c.conn.Close()
}

func (c *MQTTConn) isWebsocket() bool {
	return c.websocketConn != nil
}

func (c *MQTTConn) nextReader() (*bufio.Reader, error) {
	if c.websocketConn != nil {
		_, r, err := c.websocketConn.NextReader()
		if err != nil {
			return nil, err
		}
		return bufio.NewReader(r), nil
	}
	return bufio.NewReader(c.conn), nil
}

func (c *MQTTConn) Write(b []byte) error {
	if c.websocketConn != nil {
		return c.websocketConn.WriteMessage(websocket.BinaryMessage, b)
	}
	_, err := c.conn.Write(b)
	return err
}
