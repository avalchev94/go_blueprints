package main

import (
	"github.com/gorilla/websocket"
)

// single chatting user
type client struct {
	// the socket for this client
	socket *websocket.Conn
	// send is a channel on which messages are sent
	send chan []byte
	// the room this client is using
	room *room
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
