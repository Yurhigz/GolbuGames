package ws

// Le fonctionnement avec un système de hubmanager va permettre de créer des rooms de communication.
// A partir du moment où un client ouvre une ws avec le serveur alors on va l'associer
// à une room, et on l'associera à la même room que son adversaire
// On crée un hubmanager qui n'est ni plus ni moins qu'une liste des rooms

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

func (hub *Hub) Run() {

}
