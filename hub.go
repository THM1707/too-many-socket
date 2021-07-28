package main

import (
	"fmt"
)

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			fmt.Printf("(!) There are %v client(s) in the room\n", len(h.clients))
			client.conn.WriteJSON(Message{Type: 0, Id: client.id, Content: client.id})
			for c := range h.clients {
				if c != client {
					c.conn.WriteJSON(Message{Type: 2, Id: c.id, Content: "connect"})
				}
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			for c := range h.clients {
				c.conn.WriteJSON(Message{Type: 2, Id: c.id, Content: "disconnect"})
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
