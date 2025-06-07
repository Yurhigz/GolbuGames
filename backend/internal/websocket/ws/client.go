package ws

import "net"

type Client struct {
	clientId string
	Conn     net.Conn
}

// Représente une connexion WebSocket
// Gère la lecture/écriture des messages
// Maintient la connexion active ping récurrent
// Gère la déconnexion
