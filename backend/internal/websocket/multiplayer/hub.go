package multiplayer

import (
	"context"
	"fmt"
	"golbugames/backend/internal/sudoku/repository"
	"golbugames/backend/internal/websocket"
	"golbugames/backend/internal/websocket/protocol"
	"log"
	"runtime/debug"
	"sync"
	"time"
)

const (
	gameWaiting  = 0
	gamesOngoing = 1
	gameFinished = 2
	gameAborted  = 3
	gamePaused   = 4
	winP1        = 0
	draw         = 1
	winP2        = 2
)

type GameStatus int

type HubManager struct {
	hubs        map[string]*Hub
	mu          sync.Mutex
	ClientQueue []*Client
}

type Hub struct {
	hubManager     *HubManager
	gameState      int
	clients        [2]*Client
	register       chan *Client
	unregister     chan *Client
	broadcast      chan websocket.Frame
	hubId          string
	playable       []int
	solution       []int
	mu             sync.RWMutex // Protection pour les opérations sur les clients
	running        bool         // Flag pour indiquer si le hub est en cours d'exécution
	completionTime int
	ctx            context.Context
}

func NewHubManager() *HubManager {
	return &HubManager{
		hubs: make(map[string]*Hub),
	}
}

func (hm *HubManager) CreateHub(ctx context.Context, matchId string) (*Hub, error) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	hub := &Hub{
		hubManager: hm,
		broadcast:  make(chan websocket.Frame, 10), // Buffer pour éviter les blocages
		register:   make(chan *Client, 10),
		unregister: make(chan *Client, 10),
		clients:    [2]*Client{nil, nil},
		running:    true,
		gameState:  gameWaiting,
		ctx:        ctx,
	}

	grid, err := repository.GetRandomDifficultyGridDB(ctx)
	if err != nil {
		log.Printf("Error retrieving random grid : %v", err)
		return nil, err
	}
	hub.playable = grid.Board
	hub.solution = grid.Solution

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

func (h *Hub) GetOpponent(client *Client) *Client {
	if h.clients[0] == client {
		return h.clients[1]
	}
	return h.clients[0]
}

func (h *Hub) HandleVictory(winner *Client, reason int) {
	loser := h.GetOpponent(winner)

	h.ProcessGameEnd(winner, loser, reason)

}

// Handler métiers pour les différentes manières de finir une partie : victoire, déconnexion ...
func (h *Hub) ProcessGameEnd(winner, loser *Client, reason int) {

	switch reason {
	case protocol.PlayerDisconnected:
		if winner != nil && !winner.closed {
			winner.SendLeavingOpponent()
		}
	case protocol.SystemGameOver:
		if loser != nil && !loser.closed {
			loser.SendGameOver(winner.baseClient.ClientId)
		}
	}

	go h.SaveGameResult(winner, loser)
	// Une fois les messages envoyés et les sauvegardes effectuées on clean et supprime le hub
	if winner != nil && !winner.closed {
		winner.cleanup()
	}
	if loser != nil && !loser.closed {
		loser.cleanup()
	}
	h.hubManager.RemoveHub(h.hubId)

}

// Logique métier de victoire
// Gérer la partie contexte + structure avec timestamp
func (h *Hub) SaveGameResult(winner, loser *Client) {
	ctx := context.Background()

	completionTime := h.completionTime

	var winnerID, loserID string
	var result int

	// Gérer le getUserID avec la synchronisation DB et l'authentification

	if winner == h.clients[0] {
		result = winP1
		winnerID = winner.baseClient.ClientId
		loserID = loser.baseClient.ClientId
	} else {
		result = winP2
		winnerID = loser.baseClient.ClientId
		loserID = winner.baseClient.ClientId
	}
	fmt.Printf("winnerId : %v \n loserId : %v \n result : %v \n completionTime: %v /n", winnerID, loserID, result, completionTime)
	if err := repository.SubmitMultiGameDB(ctx, winnerID, loserID, result, completionTime); err != nil {
		log.Printf("Error saving game result: %v", err)
	}
}

// Hub Running
func (h *Hub) run() {
	log.Printf("Hub %s started", h.hubId)

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Hub %s panic recovered: %v\n%s", h.hubId, r, string(debug.Stack()))
			// log.Printf("Hub %s panic recovered: %v", h.hubId, r)
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
			opponent := h.GetOpponent(client)
			if h.gameState == gamesOngoing && opponent != nil {
				log.Printf("Game was running, player %v left, ending the game", client.baseClient.ClientId)
				h.mu.Lock()
				h.completionTime = 0
				h.mu.Unlock()
				h.HandleVictory(opponent, protocol.PlayerLeft)
			} else {
				log.Printf("Hub was waiting for another player but player %v left", client.baseClient.ClientId)
				h.handleClientUnregister(client)
			}

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
		client.hub = h
		h.clients[0].SendWaitingOpponent()

	} else if h.clients[1] == nil {
		h.clients[1] = client
		client.hub = h

		for _, c := range h.clients {
			if c != nil {
				c.SendGameStart()
			}
		}

		// Changer l'état du jeu
		h.gameState = gamesOngoing
	}
}

func (h *Hub) handleClientUnregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[0] == client {
		h.clients[0] = nil
		if h.clients[1] != nil {
			h.clients[1].SendWaitingOpponent()
		}
		log.Printf("Client unregistered from hub %s (was player 1)", h.hubId)

	} else if h.clients[1] == client {
		h.clients[1] = nil
		if h.clients[0] != nil {
			h.clients[0].SendWaitingOpponent()
		}
		log.Printf("Client unregistered from hub %s (was player 2)", h.hubId)
	}

	// Si plus aucun client, marquer pour suppression
	if h.clientCount() == 0 {
		h.gameState = gameAborted
		// Contrôler le comportement dans ce cas
		h.running = false
		log.Printf("Hub %s is empty, marking for cleanup", h.hubId)
	}
}

func (h *Hub) handleBroadcast(message websocket.Frame) {
	h.mu.RLock()
	defer h.mu.RUnlock()

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
