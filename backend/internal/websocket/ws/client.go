package ws

import (
	"net"
	"sync"
)

type Client struct {
	clientId string
	Conn     net.Conn
	mu       sync.Mutex
	send     chan []byte
	hub      *Hub
	matchId  string
}

// Représente une connexion WebSocket
// Gère la lecture/écriture des messages
// Maintient la connexion active ping récurrent
// Gère la déconnexion

func newClient() *Client {
	return &Client{}
}

func (c *Client) writePump() {

}

func (c *Client) readPump() {

}

//  Implémenter l'attribution du matchmaking
//  l'attribution des ID et la création des hubs en conséquence
