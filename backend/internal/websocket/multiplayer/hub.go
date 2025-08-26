package multiplayer

// Le fonctionnement avec un système de hubmanager va permettre de créer des rooms de communication.
// A partir du moment où un client ouvre une ws avec le serveur alors on va l'associer
// à une room, et on l'associera à la même room que son adversaire
// On crée un hubmanager qui n'est ni plus ni moins qu'une liste des rooms

type HubManager struct {
	hubs map[string]*Hub
}

type Hub struct {
	clients    [2]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	hubId      string
}

func NewHubManager() *HubManager {
	return &HubManager{
		hubs: make(map[string]*Hub),
	}
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    [2]*Client{nil, nil},
	}
}

func (hm *HubManager) CreateHub(matchId string) *Hub {
	hub := newHub()
	hub.hubId = matchId
	hm.hubs[matchId] = hub
	go hub.run()
	return hub
}

func (hm *HubManager) GetHub(matchId string) *Hub {
	return hm.hubs[matchId]
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			if h.clients[0] == nil {
				h.clients[0] = client
				client.send <- []byte("Waiting for opponent...")
			} else if h.clients[1] == nil {
				h.clients[1] = client

				message := []byte("Opponent found... Game starting!")
				h.clients[0].send <- message
				h.clients[1].send <- message
			}
		case client := <-h.unregister:
			if h.clients[0] == client {
				h.clients[0] = nil
				if h.clients[1] != nil {
					h.clients[1].send <- []byte("Opponent disconnected")
				}
			} else if h.clients[1] == client {
				h.clients[1] = nil
				if h.clients[0] != nil {
					h.clients[0].send <- []byte("Opponent disconnected")
				}
			}

		case message := <-h.broadcast:
			if h.clients[0] != nil {
				h.clients[0].send <- message
			}
			if h.clients[1] != nil {
				h.clients[1].send <- message
			}
		}
	}
}
