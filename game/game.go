package game

import (
    "errors"
    "math/rand"
)

// Position represents a cell position in the game grid
type Position struct {
    Row    int
    Column int
}

// Player represents a single player's state
type Player struct {
    // Grid stores the dice values. 0 represents an empty cell
    Grid      [3][3]int
    // Score is calculated based on dice combinations
    Score     int
    // Name helps identify players in logs and UI
    Name      string
}

// GameState represents the complete state of a Knucklebones game
type GameState struct {
    // Players array holds exactly two players
    Players        [2]Player
    // CurrentPlayer is the index (0 or 1) of the active player
    CurrentPlayer  int
    // CurrentDie is the die value that must be placed
    CurrentDie     int
    // GameOver indicates if the game has concluded
    GameOver      bool
    // Winner stores the index of the winning player (-1 if game isn't over)
    Winner        int
}

// NewGame initializes a fresh game state
func NewGame(player1Name, player2Name string) *GameState {
    return &GameState{
        Players: [2]Player{
            {Name: player1Name, Grid: [3][3]int{}},
            {Name: player2Name, Grid: [3][3]int{}},
        },
        CurrentPlayer: 0,
        CurrentDie:    rollDie(),
        Winner:        -1,
    }
}

// rollDie simulates rolling a 6-sided die
func rollDie() int {
    return rand.Intn(6) + 1
}

// IsValidMove checks if a move is legal
func (g *GameState) IsValidMove(column int) bool {
    if column < 0 || column >= 3 {
        return false
    }

    // Find the first empty row in the column
    for row := 0; row < 3; row++ {
        if g.Players[g.CurrentPlayer].Grid[row][column] == 0 {
            return true
        }
    }
    return false
}

// MakeMove executes a player's move and updates the game state
func (g *GameState) MakeMove(column int) error {
    if g.GameOver {
        return errors.New("game is already over")
    }

    if !g.IsValidMove(column) {
        return errors.New("invalid move")
    }

    // Find the first empty row and place the die
    for row := 0; row < 3; row++ {
        if g.Players[g.CurrentPlayer].Grid[row][column] == 0 {
            g.Players[g.CurrentPlayer].Grid[row][column] = g.CurrentDie
            break
        }
    }

    // Remove matching dice from opponent's column
    opponentIdx := (g.CurrentPlayer + 1) % 2
    for row := 0; row < 3; row++ {
        if g.Players[opponentIdx].Grid[row][column] == g.CurrentDie {
            g.Players[opponentIdx].Grid[row][column] = 0
        }
    }

    // Update scores for both players
    g.updateScores()

    // Check if game is over
    if g.isGameOver() {
        g.GameOver = true
        g.determineWinner()
    } else {
        // Prepare for next turn
        g.CurrentPlayer = opponentIdx
        g.CurrentDie = rollDie()
    }

    return nil
}

// updateScores recalculates scores for both players
func (g *GameState) updateScores() {
    for i := range g.Players {
        score := 0
        // Calculate score for each column
        for col := 0; col < 3; col++ {
            score += g.calculateColumnScore(i, col)
        }
        g.Players[i].Score = score
    }
}

// calculateColumnScore computes the score for a single column
func (g *GameState) calculateColumnScore(playerIdx, column int) int {
    // Count occurrences of each die value
    counts := make(map[int]int)
    for row := 0; row < 3; row++ {
        value := g.Players[playerIdx].Grid[row][column]
        if value != 0 {
            counts[value]++
        }
    }

    // Calculate score based on combinations
    score := 0
    for value, count := range counts {
        // Multiply the die value by its count squared
        // This gives bonus points for matching dice
        score += value * count * count
    }
    return score
}

// isGameOver checks if either player has a full grid by examining each player's grid
// in turn. Uses goto for efficient early exit from nested loops - this is one of the
// few cases where goto improves code clarity and performance. As soon as we find an
// empty cell in a player's grid, we can skip checking their remaining cells and move
// to the next player.
func (g *GameState) isGameOver() bool {
    for playerIdx := range g.Players {
        for row := 0; row < 3; row++ {
            for col := 0; col < 3; col++ {
                // If we find any empty cell, this player's grid isn't full
                if g.Players[playerIdx].Grid[row][col] == 0 {
                    // Continue checking next player
                    goto nextPlayer
                }
            }
        }
        // If we get here, no empty cells were found - game is over
        return true
    nextPlayer:
    }
    return false
}

// determineWinner sets the winner based on final scores
func (g *GameState) determineWinner() {
    if g.Players[0].Score > g.Players[1].Score {
        g.Winner = 0
    } else if g.Players[1].Score > g.Players[0].Score {
        g.Winner = 1
    } else {
        g.Winner = -1 // Draw
    }
}
