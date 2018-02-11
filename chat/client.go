package main

import (
	"github.com/gorilla/websocket"
	"time"
)

// single chatting user
type client struct {
	// the socket for this client
	socket *websocket.Conn
	// send is a channel on which messages are sent
	send chan *message
	// the room this client is using
	room *room
	// information about the user
	userData map[string]interface{}
}

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
