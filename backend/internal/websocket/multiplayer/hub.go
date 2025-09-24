package multiplayer

import (
	"context"
	"fmt"
	"golbugames/backend/internal/websocket"
	"golbugames/backend/internal/websocket/protocol"
	"log"
	"sync"
	"time"
)

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
	mu          sync.Mutex
	ClientQueue []*Client
}

type Hub struct {
	gameState  int
	clients    [2]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan websocket.Frame
	hubId      string
	playable   []int
	solution   []int
	mu         sync.RWMutex // Protection pour les opérations sur les clients
	running    bool         // Flag pour indiquer si le hub est en cours d'exécution
}

func NewHubManager() *HubManager {
	return &HubManager{
		hubs: make(map[string]*Hub),
	}
}

func newHub(ctx context.Context) (*Hub, error) {
	return &Hub{
		broadcast:  make(chan websocket.Frame, 10), // Buffer pour éviter les blocages
		register:   make(chan *Client, 10),
		unregister: make(chan *Client, 10),
		clients:    [2]*Client{nil, nil},
		running:    true,
		gameState:  gameWaiting,
	}, nil
}

func (hm *HubManager) CreateHub(ctx context.Context, matchId string) (*Hub, error) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	hub, err := newHub(ctx)
	if err != nil {
		return nil, err
	}
	hub.hubId = matchId
	hm.hubs[matchId] = hub
	go hub.run()
	return hub, nil
}

func (hm *HubManager) GetHub(matchId string) *Hub {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	return hm.hubs[matchId]
}

func (hm *HubManager) RemoveHub(matchId string) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	if hub, exists := hm.hubs[matchId]; exists {
		hub.mu.Lock()
		hub.running = false
		hub.mu.Unlock()

		// Fermer les canaux de manière sécurisée
		go func() {
			time.Sleep(100 * time.Millisecond)
			close(hub.register)
			close(hub.unregister)
			close(hub.broadcast)
		}()

		delete(hm.hubs, matchId)
		log.Printf("Hub %s removed", matchId)
	}
}

func (h *Hub) run() {
	log.Printf("Hub %s started", h.hubId)

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Hub %s panic recovered: %v", h.hubId, r)
		}
		log.Printf("Hub %s stopped", h.hubId)
	}()

	for h.running {
		select {
		case client, ok := <-h.register:
			if !ok {
				log.Printf("Register channel closed for hub %s", h.hubId)
				return
			}
			h.handleClientRegister(client)

		case client, ok := <-h.unregister:
			if !ok {
				log.Printf("Unregister channel closed for hub %s", h.hubId)
				return
			}
			h.handleClientUnregister(client)

		case message, ok := <-h.broadcast:
			if !ok {
				log.Printf("Broadcast channel closed for hub %s", h.hubId)
				return
			}
			h.handleBroadcast(message)
		}
	}
}

func (h *Hub) handleClientRegister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[0] == nil {
		h.clients[0] = client
		client.hub = h                                                          // Assigner le hub au client
		payload, _ := protocol.NewSystemMessage("Waiting for opponent...", 200) // il est dans la file d'attente côté frontend
		resp := &websocket.Frame{
			Opcode:  websocket.OpcodeText,
			FIN:     true,
			Payload: payload,
		}
		log.Printf("Client %s registered to hub %s as player 1", client.baseClient.ClientId, h.hubId)

		// Envoyer de manière sécurisée
		select {
		case client.baseClient.Send <- resp:
		default:
			log.Printf("Failed to send message to client %s - channel full", client.baseClient.ClientId)
		}

	} else if h.clients[1] == nil {
		h.clients[1] = client
		client.hub = h // Assigner le hub au client
		payload, _ := protocol.NewSystemMessage("Opponent found... Game will start", 200)
		resp := &websocket.Frame{
			Opcode:  websocket.OpcodeText,
			FIN:     true,
			Payload: payload,
		}

		log.Printf("Client %s registered to hub %s as player 2", client.baseClient.ClientId, h.hubId)

		// Envoyer aux deux clients
		for i, c := range h.clients {
			if c != nil {
				select {
				case c.baseClient.Send <- resp:
					log.Printf("Sent game start message to client %d", i)
				default:
					log.Printf("Failed to send message to client %d - channel full", i)
				}
			}
		}

		// Changer l'état du jeu
		h.gameState = gamesOngoing
	}
}

func (h *Hub) handleClientUnregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	payload, _ := protocol.NewSystemMessage("Opponent disconnected", protocol.PlayerLeft)
	resp := &websocket.Frame{
		Opcode:  websocket.OpcodeText,
		FIN:     true,
		Payload: payload,
	}

	if h.clients[0] == client {
		h.clients[0] = nil
		if h.clients[1] != nil {
			select {
			case h.clients[1].baseClient.Send <- resp:
				log.Printf("Notified player 2 about player 1 disconnect")
			default:
				log.Printf("Failed to notify player 2 about disconnect")
			}
		}
		log.Printf("Client unregistered from hub %s (was player 1)", h.hubId)

	} else if h.clients[1] == client {
		h.clients[1] = nil
		if h.clients[0] != nil {
			select {
			case h.clients[0].baseClient.Send <- resp:
				log.Printf("Notified player 1 about player 2 disconnect")
			default:
				log.Printf("Failed to notify player 1 about disconnect")
			}
		}
		log.Printf("Client unregistered from hub %s (was player 2)", h.hubId)
	}

	// Si plus aucun client, marquer pour suppression
	if h.clientCount() == 0 {
		h.gameState = gameAborted
		log.Printf("Hub %s is empty, marking for cleanup", h.hubId)
	}
}

func (h *Hub) handleBroadcast(message websocket.Frame) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// payload, err := json.Marshal(respMsg)
	// if err != nil {
	// 	log.Printf("Error marshaling broadcast message: %v", err)
	// 	return
	// }

	// resp := &websocket.Frame{
	// 	Opcode:  websocket.OpcodeText,
	// 	FIN:     true,
	// 	Payload: payload,
	// }

	// Diffuser à tous les clients connectés
	for i, c := range h.clients {
		if c != nil {
			select {
			case c.baseClient.Send <- &message:
				log.Printf("Broadcasted message to client %d", i)
			default:
				log.Printf("Failed to broadcast to client %d - channel full", i)
			}
		}
	}
}

func (hm *HubManager) RemoveClientFromQueue(client *Client) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	for i, c := range hm.ClientQueue {
		if c == client {
			hm.ClientQueue = append(hm.ClientQueue[:i], hm.ClientQueue[i+1:]...)
			log.Printf("Client %s removed from matchmaking queue", client.baseClient.ClientId)
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

// handler matchmaking
func (hm *HubManager) MatchmakingLoop(ctx context.Context) {
	log.Printf("Starting matchmaking loop...")
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Matchmaking loop stopped due to context cancellation")
			return
		case <-ticker.C:
			hm.processMatchmaking(ctx)
		}
	}
}

// fonction métier pour le matchamking
func (hm *HubManager) processMatchmaking(ctx context.Context) {
	hm.mu.Lock()
	queueLength := len(hm.ClientQueue)
	hubsCount := len(hm.hubs)

	if queueLength == 0 {
		hm.mu.Unlock()
		return
	}

	queue := make([]*Client, queueLength)
	copy(queue, hm.ClientQueue)
	hm.mu.Unlock()

	log.Printf("Processing matchmaking - Queue: %d, Hubs: %d", queueLength, hubsCount)

	// Chercher un hub disponible
	var availableHub *Hub
	hm.mu.Lock()
	for _, hub := range hm.hubs {
		hub.mu.RLock()
		if hub.gameState == gameWaiting && hub.clientCount() < 2 {
			availableHub = hub
			hub.mu.RUnlock()
			break
		}
		hub.mu.RUnlock()
	}
	hm.mu.Unlock()

	// Si pas de hub disponible, en créer un
	if availableHub == nil && queueLength > 0 {
		log.Printf("No available hub found, creating a new one")
		var err error
		availableHub, err = hm.CreateHub(ctx, createId())
		if err != nil {
			log.Printf("Error creating new hub: %v", err)
			return
		}
	}

	// Assigner les clients au hub
	if availableHub != nil {
		clientsToAssign := make([]*Client, 0, 2)

		hm.mu.Lock()
		for _, client := range hm.ClientQueue {
			if client.hub == nil && len(clientsToAssign) < 2-availableHub.clientCount() {
				clientsToAssign = append(clientsToAssign, client)
			}
		}
		hm.mu.Unlock()

		for _, client := range clientsToAssign {
			select {
			case availableHub.register <- client:
				hm.RemoveClientFromQueue(client)
				log.Printf("Client %s assigned to hub %s", client.baseClient.ClientId, availableHub.hubId)
			default:
				log.Printf("Failed to register client %s to hub - channel full", client.baseClient.ClientId)
			}
		}
	}
}
