package main

type room struct {
	// that's a channel that holds incoming messages that should be
	// forwarded to the other clients
	forward chan []byte
	// channel for clients wishing to join the room
	join chan *client
	// channel for clients wishing to leave the room
	leave chan *client
	// holds all current clients in this room
	clients map[*client]bool
}
