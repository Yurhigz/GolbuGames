package ws

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// HubManager gère les hubs de communication entre les clients.
// Il permet de créer des hubs, d'enregistrer des clients, de les associer à des hubs,
// et de gérer le matchmaking entre les clients.
// Il maintient une liste de hubs et une file d'attente de clients en attente de matchmaking.
// Les clients sont associés à des hubs en fonction de leur disponibilité et de l'état du jeu.
// Chaque hub gère deux clients, leur envoi de messages et leur état de jeu.

const (
	gameWaiting  = 0 // En attente de joueur avec au moins un joueur
	gamesOngoing = 1 // Partie en cours donc joueurs = 2
	gameFinished = 2 // Partie terminée => joueurs = 0
	gameAborted  = 3 // Deconnexion ou autre => joueurs >= 1
	gamePaused   = 4 // Partie en cours donc joueurs = 2
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
	ready      chan struct{}
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
		//tmp test
		ready: make(chan struct{}),
	}
}

func (hm *HubManager) CreateHub(matchId string) *Hub {
	hub := newHub()
	hub.hubId = matchId
	hm.hubs[matchId] = hub
	go hub.run()
	//tmp test
	<-hub.ready
	log.Printf("[INFO] New hub created with ID: %s\n", matchId)
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

func createId() string {
	return fmt.Sprintf("hub_%d", time.Now().UnixNano())
}

func (h *Hub) run() {
	log.Printf("[INFO] Hub [%s] run() started\n", h.hubId)
	close(h.ready)
	for {
		select {
		case client := <-h.register:
			log.Printf("[INFO] Hub %s received client in register\n", h.hubId)
			if h.clients[0] == nil {
				h.clients[0] = client
				client.hub = h
				client.matchId = h.hubId
				time.Sleep(1000 * time.Millisecond)
				client.send <- NewTextFrame("Waiting for opponent...")
			} else {
				h.clients[1] = client
				client.hub = h
				client.matchId = h.hubId
				time.Sleep(100 * time.Millisecond)
				message := "Opponent found... Game starting!"
				h.clients[0].send <- NewTextFrame(message)
				h.clients[1].send <- NewTextFrame(message)
			}

		case client := <-h.unregister:
			log.Printf("[INFO] Hub %s received client in unregister\n", h.hubId)
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
		case message := <-h.broadcast:
			log.Printf("[INFO] Hub %s received message in broadcast\n", h.hubId)
			if h.clients[0] != nil {
				h.clients[0].send <- message
			}
			if h.clients[1] != nil {
				h.clients[1].send <- message
			}
		}
		if h.clients[0] == nil && h.clients[1] == nil {
			h.gameState = gameFinished
			return
		}
	}
}

// refactoriser avec deux fonctions pour les deux cas d'usage
func (hm *HubManager) MatchmakingLoop() {
	for {
		time.Sleep(5 * time.Second) // ajuster la fréquence de vérification si nécessaire mais dans le cadre du test/debug on peut rester à 5s
		hm.mu.Lock()
		queue := make([]*Client, len(hm.ClientQueue))
		copy(queue, hm.ClientQueue)
		hm.mu.Unlock()

		if len(hm.hubs) > 0 && len(queue) > 0 {
			log.Printf("[INFO] Current matchmaking queue length: %d\n", len(queue))
			log.Printf("[INFO] Current hubs count: %d\n", len(hm.hubs))
			availability := false
			var availableHub *Hub
			for _, hub := range hm.hubs {
				if hub.gameState == gameWaiting && hub.clientCount() < 2 {
					availability = true
					availableHub = hub
					break
				}
			}
			if !availability {
				log.Printf("[INFO] No available hub found, creating a new one\n")
				availableHub = hm.CreateHub(createId())
			}

			for _, client := range hm.ClientQueue {
				if client.matchId == "" {
					availableHub.register <- client
					hm.RemoveClientFromQueue(client)
					client.send <- NewTextFrame("You have been matched with an opponent!")
					log.Printf("[INFO] A match has been made\n")

				}
				if availableHub.clientCount() == 2 {
					log.Printf("[INFO] Hub %s is now full with 2 clients\n", availableHub.hubId)
					availableHub.gameState = gamesOngoing
					break
				}
			}

		} else if len(queue) >= 1 {
			var longestWaitingTime *Client
			for _, client := range queue {
				waiting_time := time.Since(client.queueTime)
				if longestWaitingTime == nil || waiting_time > time.Since(longestWaitingTime.queueTime) {
					longestWaitingTime = client
				}
			}
			if longestWaitingTime != nil {
				hub := hm.CreateHub(createId())
				time.Sleep(1000 * time.Millisecond)
				hub.register <- longestWaitingTime
				hm.RemoveClientFromQueue(longestWaitingTime)
			}

		}
	}

}

func (hm *HubManager) HubCleanupLoop() {
	for {
		time.Sleep(30 * time.Second) // fréquence de nettoyage à faire varier si on veut nettoyer plus ou moins vite et selon le volume de joueurs
		var toDelete []string
		hm.mu.Lock()
		for id, hub := range hm.hubs {
			if hub.clientCount() == 0 {
				toDelete = append(toDelete, id)
			}
		}
		hm.mu.Unlock()

		for _, id := range toDelete {
			hm.RemoveHub(id)
		}
	}

}
