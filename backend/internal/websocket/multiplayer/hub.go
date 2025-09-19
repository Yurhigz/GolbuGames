package multiplayer

import (
	"context"
	"encoding/json"
	"golbugames/backend/internal/sudoku/repository"
	"golbugames/backend/internal/websocket"
	"golbugames/backend/internal/websocket/protocol"
	"sync"
)

// Le fonctionnement avec un système de hubmanager va permettre de créer des rooms de communication.
// A partir du moment où un client ouvre une ws avec le serveur alors on va l'associer
// à une room, et on l'associera à la même room que son adversaire
// On crée un hubmanager qui n'est ni plus ni moins qu'une liste des rooms

type HubManager struct {
	hubs        map[string]*Hub
	mu          sync.Mutex
	ClientQueue []*Client
}

type Hub struct {
	clients    [2]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan websocket.BroadcastFrame
	hubId      string
	playable   []int
	solution   []int
}

func NewHubManager() *HubManager {
	return &HubManager{
		hubs: make(map[string]*Hub),
	}
}

func newHub(ctx context.Context) (*Hub, error) {
	grid, err := repository.GetRandomDifficultyGridDB(ctx)
	if err != nil {
		return nil, err
	}

	return &Hub{
		broadcast:  make(chan websocket.BroadcastFrame),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    [2]*Client{nil, nil},
		playable:   grid.Board,
		solution:   grid.Solution,
	}, nil
}

func (hm *HubManager) CreateHub(ctx context.Context, matchId string) (*Hub, error) {
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
	return hm.hubs[matchId]
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			if h.clients[0] == nil {
				h.clients[0] = client
				payload, _ := protocol.NewSystemMessage("Waiting for opponent...", 200)
				resp := &websocket.Frame{
					Opcode:  websocket.OpcodeText,
					FIN:     true,
					Payload: payload,
				}
				client.baseClient.Send <- resp
			} else if h.clients[1] == nil {
				h.clients[1] = client
				payload, _ := protocol.NewSystemMessage("Opponent found... Game will start", 200)
				resp := &websocket.Frame{
					Opcode:  websocket.OpcodeText,
					FIN:     true,
					Payload: payload,
				}
				h.clients[0].baseClient.Send <- resp
				h.clients[1].baseClient.Send <- resp
			}
		case client := <-h.unregister:
			payload, _ := protocol.NewSystemMessage("Opponent disconnected", protocol.PlayerLeft)
			resp := &websocket.Frame{
				Opcode:  websocket.OpcodeText,
				FIN:     true,
				Payload: payload,
			}
			if h.clients[0] == client {
				h.clients[0] = nil
				if h.clients[1] != nil {
					h.clients[1].baseClient.Send <- resp
				}
			} else if h.clients[1] == client {
				h.clients[1] = nil
				if h.clients[0] != nil {
					h.clients[0].baseClient.Send <- resp
				}
			}

		case message := <-h.broadcast:
			// Ici je dois renvoyer une frame enveloppée dans un broadcastframe à tous les clients du hub
			respMsg := protocol.ChatMessage{
				Type:    protocol.MessageTypeChat,
				Sender:  message.Sender,
				Message: string(message.Frame.Payload),
			}
			// A vérifier et tester /
			payload, _ := json.Marshal(respMsg)
			resp := &websocket.Frame{
				Opcode:  websocket.OpcodeText,
				FIN:     true,
				Payload: payload,
			}
			if h.clients[0] != nil {
				h.clients[0].baseClient.Send <- resp
			}
			if h.clients[1] != nil {
				h.clients[1].baseClient.Send <- resp
			}
		}
	}
}
