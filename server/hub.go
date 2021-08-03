package server

// Hub maintains active connections and broadcast messages
// to clients
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Incoming messages from the clients
	broadcast chan []byte

	// Register client's requests
	register chan *Client

	// Unregister requests from the clients
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub {
		clients: make(map[*Client]bool),
		broadcast:   make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				close(client.send)
				delete(h.clients, client)
			}

		case msg := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- msg:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
