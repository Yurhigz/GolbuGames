package ws

import (
	"fmt"
	"sync"
	"time"
)

// Le fonctionnement avec un système de hubmanager va permettre de créer des rooms de communication.
// A partir du moment où un client ouvre une ws avec le serveur alors on va l'associer
// à une room, et on l'associera à la même room que son adversaire
// On crée un hubmanager qui n'est ni plus ni moins qu'une liste des rooms
const (
	gameWaiting  = 0
	gamesOngoing = 1
	gameFinished = 2
	gameAborted  = 3
	gamePaused   = 4
)

type GameStatus int

type HubManager struct {
	hubs        map[string]*Hub
	ClientQueue []*Client
	mu          sync.Mutex
}

// Pas d'intégration de la différence d'élo entre les joueurs pour le moment,
// on se concentre sur la gestion des rooms et des clients avec un matchmaking basique
type Hub struct {
	clients    [2]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Frame
	hubId      string
	gameState  GameStatus
}

func NewHubManager() *HubManager {
	return &HubManager{
		hubs: make(map[string]*Hub),
	}
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Frame),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    [2]*Client{nil, nil},
	}
}

// checker utilité fonction par rapport à newHub
func (hm *HubManager) CreateHub(matchId string) *Hub {
	hub := newHub()
	hub.hubId = matchId
	hm.hubs[matchId] = hub
	go hub.run()
	fmt.Printf("[DEBUG] New hub created with ID: %s\n", matchId)
	return hub
}

func (hm *HubManager) GetHub(matchId string) *Hub {
	return hm.hubs[matchId]
}

func (hm *HubManager) RemoveHub(matchId string) {
	if hub, exists := hm.hubs[matchId]; exists {
		close(hub.register)
		close(hub.unregister)
		close(hub.broadcast)
		delete(hm.hubs, matchId)
	}
}

func (h *Hub) run() {
	fmt.Printf("[DEBUG] Hub %s run() started\n", h.hubId)
	for {
		fmt.Printf("[DEBUG] Hub %s about to enter select\n", h.hubId)
		select {

		case client := <-h.register:

			if h.clients[0] == nil {
				h.clients[0] = client
				client.hub = h
				time.Sleep(1000 * time.Millisecond)
				fmt.Printf("<hub run> client waiting for opponent")
				client.send <- NewTextFrame("Waiting for opponent...")
				fmt.Printf("<hub run> client connected to game room")
			} else if h.clients[1] == nil {
				h.clients[1] = client
				client.hub = h
				fmt.Printf("<hub run> client connected to game room")
				time.Sleep(100 * time.Millisecond)
				message := "Opponent found... Game starting!"
				h.clients[0].send <- NewTextFrame(message)
				h.clients[1].send <- NewTextFrame(message)
			}

		case client := <-h.unregister:
			fmt.Println("Processing unregister...")
			if h.clients[0] == client {
				h.clients[0] = nil
				if h.clients[1] != nil {
					h.clients[1].send <- NewTextFrame("Opponent disconnected")
				}
			} else if h.clients[1] == client {
				h.clients[1] = nil
				if h.clients[0] != nil {
					h.clients[0].send <- NewTextFrame("Opponent disconnected")
				}
			}
		// vérifier pour le broadcast si les clients envoies des messages chiffrés
		case message := <-h.broadcast:
			fmt.Println("Processing broadcast...")
			if h.clients[0] != nil {
				h.clients[0].send <- message
			}
			if h.clients[1] != nil {
				h.clients[1].send <- message
			}
			// default:
			// 	fmt.Println("Hub register channel is blocked!")
			// 	return
		}
		if h.clients[0] == nil && h.clients[1] == nil {
			h.gameState = gameFinished
			return
		}
	}
}

func (hm *HubManager) RemoveClientFromQueue(client *Client) {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	for i, c := range hm.ClientQueue {
		if c == client {
			hm.ClientQueue = append(hm.ClientQueue[:i], hm.ClientQueue[i+1:]...)
			return
		}
	}
}

func (h *Hub) clientCount() int {
	count := 0
	for _, c := range h.clients {
		if c != nil {
			count++
		}
	}
	return count
}

// Ici l'idée c'est de faire une boucle qui tourne en fond en permanence pour match les clients
// je filtre en fonction d'abord du temps d'attente puis ensuite de l'élo des joueurs
// Il faut également que je trouve un moyen de mettre à jour soit le temps d'attente, soit j'utilise
// un calcul du durée en fonction de l'heure actuelle pour éviter de faire des appels trop fréquents

// refactoriser avec deux fonctions pour les deux cas d'usage
func (hm *HubManager) MatchmakingLoop() {
	for {
		time.Sleep(5 * time.Second) // Ajuster la fréquence de vérification si nécessaire
		hm.mu.Lock()
		queue := make([]*Client, len(hm.ClientQueue))
		copy(queue, hm.ClientQueue)
		hm.mu.Unlock()
		// vérifier si un hub est déjà en cours et qu'il y a une place pour un nouveau client ainsi que des clients
		// dans la queue
		fmt.Println("[DEBUG] Matchmaking loop tick")
		fmt.Printf("Current matchmaking queue length: %d\n", len(queue))
		fmt.Printf("Current hubs count: %d\n", len(hm.hubs))
		if len(hm.hubs) > 0 && len(queue) > 0 {
			fmt.Printf("Current matchmaking queue length: %d\n", len(queue))
			fmt.Printf("Current hubs count: %d\n", len(hm.hubs))
			for _, hub := range hm.hubs {
				fmt.Printf("[DEBUG] Checking hub %s with state %d\n", hub.hubId, hub.gameState)
				if hub.gameState == gameWaiting {
					fmt.Printf("[DEBUG] Hub %s is waiting for clients\n", hub.hubId)
					if hub.clientCount() < 2 {
						// Si le hub a moins de 2 clients, on peut essayer de les associer
						for _, client := range hm.ClientQueue {
							if client.matchId == "" { // Si le client n'est pas déjà associé à une room
								fmt.Printf("[DEBUG] Avant hub.register <- client (%s)\n", client.clientId)
								hub.register <- client
								fmt.Printf("[DEBUG] Après hub.register <- client (%s)\n", client.clientId)
								client.hub = hub
								client.matchId = hub.hubId
								hm.mu.Lock()
								hm.RemoveClientFromQueue(client)
								hm.mu.Unlock()
								fmt.Printf("[DEBUG] Client %s ajouté à la queue\n", client.clientId)
								fmt.Printf("[DEBUG] Client %s envoyé dans hub %s\n", client.clientId, hub.hubId)
								fmt.Printf("[DEBUG] hub.clientCount() = %d\n", hub.clientCount())
								fmt.Printf("[DEBUG] hub.gameState = %d\n", hub.gameState)
								client.send <- NewTextFrame("You have been matched with an opponent!")
								fmt.Printf("A match has been made\n")

							}
							if hub.clientCount() == 2 {
								fmt.Printf("Hub %s is now full with 2 clients\n", hub.hubId)
								hub.gameState = gamesOngoing
								break
							}
						}
					}
				}
			}
		} else if len(queue) >= 1 {
			fmt.Printf("[DEBUG] checking for longest waiting client\n")
			var longestWaitingTime *Client
			for _, client := range queue {
				waiting_time := time.Since(client.queueTime)
				if longestWaitingTime == nil || waiting_time > time.Since(longestWaitingTime.queueTime) {
					longestWaitingTime = client
				}
			}
			if longestWaitingTime != nil {
				longestWaitingTime.matchId = longestWaitingTime.clientId + "_room"
				hub := hm.CreateHub(longestWaitingTime.matchId)
				time.Sleep(1000 * time.Millisecond)
				hub.register <- longestWaitingTime
				longestWaitingTime.hub = hub
				fmt.Printf("[DEBUG] Client %s ajouté à la queue\n", longestWaitingTime.clientId)
				hm.mu.Lock()
				hm.RemoveClientFromQueue(longestWaitingTime)
				hm.mu.Unlock()
				fmt.Printf("[DEBUG] Client %s envoyé dans hub %s\n", longestWaitingTime.clientId, hub.hubId)
			}

		}
	}

}
