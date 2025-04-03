package server

import (
	"encoding/json"
	"fmt"
	"knucklebones/game"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// GameServer manages the state of a single active game
type GameServer struct {
	// The single game instance
	game *game.GameState
}

// NewGameServer initializes a new game server with no active game
func NewGameServer() *GameServer {
	return &GameServer{
		game: nil,
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

	// Configure routes (no game IDs needed)
	router.HandleFunc("/game", adapter.handleCreateGame).Methods("POST")
	router.HandleFunc("/game", adapter.handleGetGame).Methods("GET")
	router.HandleFunc("/game/move", adapter.handleMakeMove).Methods("POST")
	router.HandleFunc("/game/complete", adapter.handleCompleteGame).Methods("POST")

	return adapter
}

func (a *RESTAdapter) StartServer() error {
	log.Printf("Starting REST server on %s", a.server.Addr)
	return a.server.ListenAndServe()
}

func (a *RESTAdapter) Stop() error {
	return a.server.Close()
}

// handleCreateGame creates a new game
func (a *RESTAdapter) handleCreateGame(w http.ResponseWriter, r *http.Request) {
	var req CreateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// If a game already exists, return an error
	if a.games.game != nil {
		http.Error(w, "A game is already in progress", http.StatusBadRequest)
		return
	}

	// Create new game
	a.games.game = game.NewGame(req.Player1Name, req.Player2Name)

	// Return the game state
	json.NewEncoder(w).Encode(GameResponse{
		CurrentState: *a.games.game,
	})
}

// handleGetGame returns the current state of the game
func (a *RESTAdapter) handleGetGame(w http.ResponseWriter, r *http.Request) {
	if a.games.game == nil {
		http.Error(w, "No active game", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(GameResponse{
		CurrentState: *a.games.game,
	})
}

// handleMakeMove processes a move in the game
func (a *RESTAdapter) handleMakeMove(w http.ResponseWriter, r *http.Request) {
	if a.games.game == nil {
		http.Error(w, "No active game", http.StatusNotFound)
		return
	}

	var req MakeMoveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := a.games.game.MakeMove(req.Column)

	if err != nil {
		json.NewEncoder(w).Encode(GameResponse{
			CurrentState: *a.games.game,
			Error:        err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(GameResponse{
		CurrentState: *a.games.game,
	})
}

// handleCompleteGame explicitly marks the game as complete and removes it
func (a *RESTAdapter) handleCompleteGame(w http.ResponseWriter, r *http.Request) {
	if a.games.game == nil {
		http.Error(w, "No active game", http.StatusNotFound)
		return
	}

	// Clear the game
	a.games.game = nil

	log.Printf("Game marked as complete and removed")

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Game successfully completed and resources released",
	})
}

// For illustrative purposes - this should be removed in production as it's redundant
func main() {
	gameServer := NewGameServer()
	restAdapter := NewRESTAdapter(gameServer, 8080)

	log.Fatal(restAdapter.StartServer())
}
