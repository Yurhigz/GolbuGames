package ws

type HubManager struct {
	hubs map[*Hub]bool
}

type Hub struct {
	clients    [2]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	hubId      int
}

func newHub() *Hub {
	return &Hub{}
}
