# Go Knucklebones

A Go implementation of the Knucklebones dice game with a RESTful API server.

## Description

Knucklebones is a tactical dice placement game for two players. Players take turns rolling a die and placing it in their grid, with the goal of creating matching combinations for higher scores while removing their opponent's matching dice.

## Features

- Full implementation of Knucklebones game logic
- RESTful API server
- Concurrent game sessions support
- Comprehensive test coverage
- Score calculation system
- Game state management

## Technical Stack

- Go 1.23.1
- Gorilla Mux for HTTP routing

## Installation

1. Clone the repository:
```bash
git clone [repository-url]
cd go-knucklebones
```

2. Install dependencies:
```bash
go mod download
```

## Usage

### Starting the Server

Run the server with:
```bash
go run main.go
```

The server will start on port 8080 by default.

### API Endpoints

#### Create a New Game
```http
POST /games
Content-Type: application/json

{
    "player1_name": "Player1",
    "player2_name": "Player2"
}
```

#### Get Game State
```http
GET /games/{id}
```

#### Make a Move
```http
POST /games/{id}/move
Content-Type: application/json

{
    "column": 0
}
```

## Game Rules

1. Players take turns rolling a die and placing it in their 3x3 grid
2. When placing a die, if the opponent has matching dice in the same column, they are removed
3. Scoring:
   - Single die: Value × 1²
   - Two matching dice: Value × 2²
   - Three matching dice: Value × 3²
4. Game ends when either player fills their grid
5. Player with the highest score wins

## Testing

Run the test suite with:
```bash
go test ./...
```

## Project Structure

```
go-knucklebones/
├── game/
│   ├── game.go         # Core game logic
│   └── game_test.go    # Game logic tests
├── server/
│   └── server.go       # HTTP server implementation
├── main.go             # Application entry point
├── go.mod              # Go module file
└── go.sum              # Go module checksum
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request
