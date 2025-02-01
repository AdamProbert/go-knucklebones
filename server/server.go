package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"knucklebones/game"

	"github.com/gorilla/mux"
)

// GameServer manages the state of all active games
type GameServer struct {
	// games stores active games using a string ID as the key
	games map[string]*game.GameState
	// mutex prevents concurrent access to the games map
	mutex sync.RWMutex
}

// NewGameServer initializes a new game server
func NewGameServer() *GameServer {
	return &GameServer{
		games: make(map[string]*game.GameState),
	}
}

// NetworkAdapter defines the interface for different networking protocols
// This allows us to swap between REST, WebSocket, etc.
type NetworkAdapter interface {
	StartServer() error
	Stop() error
}

// RESTAdapter implements the NetworkAdapter interface for REST
type RESTAdapter struct {
	server *http.Server
	games  *GameServer
}

// Request/Response structures for our REST API
type CreateGameRequest struct {
	Player1Name string `json:"player1_name"`
	Player2Name string `json:"player2_name"`
}

type GameResponse struct {
	ID           string         `json:"id"`
	CurrentState game.GameState `json:"state"`
	Error        string         `json:"error,omitempty"`
}

type MakeMoveRequest struct {
	Column int `json:"column"`
}

// NewRESTAdapter creates a new REST adapter with routes configured
func NewRESTAdapter(gameServer *GameServer, port int) *RESTAdapter {
	router := mux.NewRouter()
	adapter := &RESTAdapter{
		games: gameServer,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: router,
		},
	}

	// Configure routes
	router.HandleFunc("/games", adapter.handleCreateGame).Methods("POST")
	router.HandleFunc("/games/{id}", adapter.handleGetGame).Methods("GET")
	router.HandleFunc("/games/{id}/move", adapter.handleMakeMove).Methods("POST")

	return adapter
}

func (a *RESTAdapter) StartServer() error {
	log.Printf("Starting REST server on %s", a.server.Addr)
	return a.server.ListenAndServe()
}

func (a *RESTAdapter) Stop() error {
	return a.server.Close()
}

// handleCreateGame creates a new game and returns its ID
func (a *RESTAdapter) handleCreateGame(w http.ResponseWriter, r *http.Request) {
	var req CreateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create new game with unique ID (using timestamp for simplicity)
	gameID := fmt.Sprintf("game_%d", time.Now().UnixNano())
	game := game.NewGame(req.Player1Name, req.Player2Name)

	// Store the game
	a.games.mutex.Lock()
	a.games.games[gameID] = game
	a.games.mutex.Unlock()

	// Return the game state
	json.NewEncoder(w).Encode(GameResponse{
		ID:           gameID,
		CurrentState: *game,
	})
}

// handleGetGame returns the current state of a game
func (a *RESTAdapter) handleGetGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["id"]

	a.games.mutex.RLock()
	game, exists := a.games.games[gameID]
	a.games.mutex.RUnlock()

	if !exists {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(GameResponse{
		ID:           gameID,
		CurrentState: *game,
	})
}

// handleMakeMove processes a move in the game
func (a *RESTAdapter) handleMakeMove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["id"]

	var req MakeMoveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	a.games.mutex.Lock()
	game, exists := a.games.games[gameID]
	if !exists {
		a.games.mutex.Unlock()
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	err := game.MakeMove(req.Column)
	a.games.mutex.Unlock()

	if err != nil {
		json.NewEncoder(w).Encode(GameResponse{
			ID:           gameID,
			CurrentState: *game,
			Error:        err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(GameResponse{
		ID:           gameID,
		CurrentState: *game,
	})
}

func main() {
	gameServer := NewGameServer()
	restAdapter := NewRESTAdapter(gameServer, 8080)

	log.Fatal(restAdapter.StartServer())
}
