package multiplayer

import (
	"context"
	"encoding/json"
	"fmt"
	"golbugames/internal/sudoku/repository"
	"golbugames/internal/websocket"
	"golbugames/internal/websocket/protocol"
	"log"
	"runtime/debug"
	"sync"
	"sync/atomic"
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
	gamesEnded  atomic.Uint64
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
	mu             sync.RWMutex
	completionTime int
	startedAt      time.Time
	ctx            context.Context
	cancel         context.CancelFunc
}

type Metrics struct {
	WaitingPlayers  int    `json:"waiting_players"`
	ActiveGames     int    `json:"active_games"`
	GamesEndedTotal uint64 `json:"games_ended_total"`
}

func NewHubManager() *HubManager {
	return &HubManager{
		hubs: make(map[string]*Hub),
	}
}

func (hm *HubManager) CreateHub(ctx context.Context, matchId string) (*Hub, error) {
	hubCtx, cancel := context.WithCancel(ctx)
	grid, err := repository.GetRandomDifficultyGridDB(hubCtx)
	if err != nil {
		cancel()
		return nil, err
	}

	hub := &Hub{
		hubManager: hm,
		broadcast:  make(chan websocket.Frame, 10),
		register:   make(chan *Client, 10),
		unregister: make(chan *Client, 10),
		gameState:  gameWaiting,
		hubId:      matchId,
		playable:   append([]int(nil), grid.Board...),
		solution:   append([]int(nil), grid.Solution...),
		ctx:        hubCtx,
		cancel:     cancel,
	}

	hm.mu.Lock()
	if _, exists := hm.hubs[matchId]; exists {
		hm.mu.Unlock()
		cancel()
		return nil, fmt.Errorf("hub %q already exists", matchId)
	}
	hm.hubs[matchId] = hub
	hm.mu.Unlock()

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
	hub, exists := hm.hubs[matchId]
	if exists {
		delete(hm.hubs, matchId)
	}
	hm.mu.Unlock()

	if exists {
		hub.cancel()
		log.Printf("Hub %s removed", matchId)
	}
}

func (h *Hub) GetOpponent(client *Client) *Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if h.clients[0] == client {
		return h.clients[1]
	}
	if h.clients[1] == client {
		return h.clients[0]
	}
	return nil
}

func (h *Hub) isOngoing() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.gameState == gamesOngoing
}

func (h *Hub) HandleMove(client *Client, move protocol.GameMessage) {
	h.mu.Lock()
	isPlayer := h.clients[0] == client || h.clients[1] == client
	valid := h.gameState == gamesOngoing &&
		isPlayer &&
		move.Position >= 0 &&
		move.Position < len(h.solution) &&
		move.Value >= 1 &&
		move.Value <= 9 &&
		move.Position < len(h.playable) &&
		h.playable[move.Position] == 0 &&
		h.solution[move.Position] == move.Value

	if valid {
		h.playable[move.Position] = move.Value
	}
	clients := h.clients
	complete := valid && h.gameState == gamesOngoing && h.isBoardComplete()
	h.mu.Unlock()

	payload, err := jsonMarshalMove(move, valid)
	if err != nil {
		log.Printf("Error marshalling move validation: %v", err)
		return
	}
	if valid {
		for _, player := range clients {
			if player != nil {
				player.SendMessage(payload)
			}
		}
	} else {
		client.SendMessage(payload)
	}

	if complete {
		h.HandleVictory(client, protocol.SystemGameOver)
	}
}

func jsonMarshalMove(move protocol.GameMessage, valid bool) ([]byte, error) {
	return json.Marshal(protocol.ValidationMove{
		Type:     protocol.MessageTypeValidateMove,
		Position: move.Position,
		Value:    move.Value,
		Valid:    valid,
	})
}

func (h *Hub) isBoardComplete() bool {
	for _, value := range h.playable {
		if value == 0 {
			return false
		}
	}
	return len(h.playable) == 81
}

func (h *Hub) HandleVictory(winner *Client, reason int) {
	h.mu.Lock()
	if h.gameState != gamesOngoing || (h.clients[0] != winner && h.clients[1] != winner) {
		h.mu.Unlock()
		return
	}
	loser := h.clients[0]
	if loser == winner {
		loser = h.clients[1]
	}
	if loser == nil {
		h.mu.Unlock()
		return
	}
	if reason == protocol.SystemGameOver && !h.isBoardComplete() {
		h.mu.Unlock()
		return
	}

	h.gameState = gameFinished
	h.completionTime = int(time.Since(h.startedAt).Seconds())
	if h.completionTime < 1 {
		h.completionTime = 1
	}
	h.hubManager.gamesEnded.Add(1)
	h.mu.Unlock()

	h.ProcessGameEnd(winner, loser, reason)
}

func (h *Hub) ProcessGameEnd(winner, loser *Client, reason int) {
	switch reason {
	case protocol.PlayerLeft, protocol.PlayerDisconnected:
		winner.SendLeavingOpponent()
	case protocol.SystemGameOver:
		winner.SendGameOver(winner.username)
		loser.SendGameOver(winner.username)
	}

	go h.SaveGameResult(winner)
	h.hubManager.RemoveHub(h.hubId)
}

func (h *Hub) SaveGameResult(winner *Client) {
	h.mu.RLock()
	player1, player2 := h.clients[0], h.clients[1]
	result := winP1
	if winner == player2 {
		result = winP2
	}
	completionTime := h.completionTime
	h.mu.RUnlock()

	if player1 == nil || player2 == nil {
		log.Printf("Cannot save result for hub %s: both players are required", h.hubId)
		return
	}
	if err := repository.SubmitMultiGameDB(
		context.Background(),
		player1.userID,
		player2.userID,
		result,
		completionTime,
	); err != nil {
		log.Printf("Error saving game result for hub %s: %v", h.hubId, err)
	}
}

func (h *Hub) run() {
	log.Printf("Hub %s started", h.hubId)
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Hub %s panic recovered: %v\n%s", h.hubId, r, string(debug.Stack()))
		}
		log.Printf("Hub %s stopped", h.hubId)
	}()

	for {
		select {
		case <-h.ctx.Done():
			return
		case client := <-h.register:
			h.handleClientRegister(client)
		case client := <-h.unregister:
			if h.isOngoing() {
				opponent := h.GetOpponent(client)
				if opponent != nil {
					h.HandleVictory(opponent, protocol.PlayerLeft)
					continue
				}
			}
			if h.handleClientUnregister(client) {
				h.hubManager.RemoveHub(h.hubId)
			}
		case message := <-h.broadcast:
			h.handleBroadcast(message)
		}
	}
}

func (h *Hub) handleClientRegister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.gameState != gameWaiting || client.isClosed() {
		return
	}
	if h.clients[0] == nil {
		if !client.setHub(h) {
			return
		}
		h.clients[0] = client
	} else if h.clients[1] == nil {
		if !client.setHub(h) {
			return
		}
		h.clients[1] = client
	} else {
		return
	}
	if h.clients[1] == nil {
		h.clients[0].SendWaitingOpponent()
		return
	}

	h.gameState = gamesOngoing
	h.startedAt = time.Now()
	for _, player := range h.clients {
		player.SendGameStart()
	}
}

func (h *Hub) handleClientUnregister(client *Client) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[0] == client {
		h.clients[0] = nil
	} else if h.clients[1] == client {
		h.clients[1] = nil
	} else {
		return h.clientCount() == 0
	}

	if h.clientCount() == 0 {
		h.gameState = gameAborted
		return true
	}
	for _, player := range h.clients {
		if player != nil {
			player.SendWaitingOpponent()
		}
	}
	return false
}

func (h *Hub) handleBroadcast(message websocket.Frame) {
	h.mu.RLock()
	clients := h.clients
	h.mu.RUnlock()
	for _, client := range clients {
		if client != nil {
			client.SendMessage(message.Payload)
		}
	}
}

func (hm *HubManager) RemoveClientFromQueue(client *Client) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	for i, queued := range hm.ClientQueue {
		if queued == client {
			hm.ClientQueue = append(hm.ClientQueue[:i], hm.ClientQueue[i+1:]...)
			log.Printf("Client %s removed from matchmaking queue", client.baseClient.ClientId)
			return
		}
	}
}

func (h *Hub) clientCount() int {
	count := 0
	for _, client := range h.clients {
		if client != nil {
			count++
		}
	}
	return count
}

func createId() string {
	return fmt.Sprintf("hub_%d", time.Now().UnixNano())
}

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

func (hm *HubManager) processMatchmaking(ctx context.Context) {
	hm.mu.Lock()
	if len(hm.ClientQueue) == 0 {
		hm.mu.Unlock()
		return
	}
	queue := append([]*Client(nil), hm.ClientQueue...)
	var availableHub *Hub
	for _, hub := range hm.hubs {
		hub.mu.RLock()
		available := hub.gameState == gameWaiting && hub.clientCount() < 2
		hub.mu.RUnlock()
		if available {
			availableHub = hub
			break
		}
	}
	hm.mu.Unlock()

	if availableHub == nil {
		var err error
		availableHub, err = hm.CreateHub(ctx, createId())
		if err != nil {
			log.Printf("Error creating multiplayer hub: %v", err)
			return
		}
	}

	availableHub.mu.RLock()
	availableSlots := 2 - availableHub.clientCount()
	availableHub.mu.RUnlock()
	if availableSlots <= 0 {
		return
	}
	if availableSlots > len(queue) {
		availableSlots = len(queue)
	}

	for _, client := range queue[:availableSlots] {
		select {
		case availableHub.register <- client:
			hm.RemoveClientFromQueue(client)
			log.Printf("Client %s assigned to hub %s", client.baseClient.ClientId, availableHub.hubId)
		case <-ctx.Done():
			return
		case <-availableHub.ctx.Done():
			return
		}
	}
}

func (hm *HubManager) MetricsSnapshot() Metrics {
	hm.mu.Lock()
	waitingPlayers := len(hm.ClientQueue)
	hubs := make([]*Hub, 0, len(hm.hubs))
	for _, hub := range hm.hubs {
		hubs = append(hubs, hub)
	}
	hm.mu.Unlock()

	activeGames := 0
	for _, hub := range hubs {
		hub.mu.RLock()
		if hub.gameState == gamesOngoing {
			activeGames++
		}
		hub.mu.RUnlock()
	}

	return Metrics{
		WaitingPlayers:  waitingPlayers,
		ActiveGames:     activeGames,
		GamesEndedTotal: hm.gamesEnded.Load(),
	}
}
