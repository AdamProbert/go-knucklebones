package game

import (
	"testing"
)

// TestNewGame verifies that a new game is initialized with correct default values
func TestNewGame(t *testing.T) {
    game := NewGame("Player1", "Player2")

    // Check initial game state
    if game.CurrentPlayer != 0 {
        t.Errorf("Expected first player to be 0, got %d", game.CurrentPlayer)
    }

    if game.GameOver {
        t.Error("New game should not be marked as over")
    }

    if game.Winner != -1 {
        t.Errorf("Expected winner to be -1, got %d", game.Winner)
    }

    // Verify player names
    if game.Players[0].Name != "Player1" {
        t.Errorf("Expected Player1, got %s", game.Players[0].Name)
    }
    if game.Players[1].Name != "Player2" {
        t.Errorf("Expected Player2, got %s", game.Players[1].Name)
    }

    // Check that initial die roll is valid
    if game.CurrentDie < 1 || game.CurrentDie > 6 {
        t.Errorf("Current die value %d is outside valid range 1-6", game.CurrentDie)
    }
}

// TestIsValidMove checks the move validation logic
func TestIsValidMove(t *testing.T) {
    game := NewGame("Player1", "Player2")

    // Test valid moves
    if !game.IsValidMove(0) {
        t.Error("Move in empty column should be valid")
    }

    // Test invalid column numbers
    if game.IsValidMove(-1) {
        t.Error("Move in column -1 should be invalid")
    }
    if game.IsValidMove(3) {
        t.Error("Move in column 3 should be invalid")
    }

    // Fill a column and verify moves become invalid
    for row := 0; row < 3; row++ {
        game.Players[0].Grid[row][0] = 1
    }
    if game.IsValidMove(0) {
        t.Error("Move in full column should be invalid")
    }
}

// TestMakeMove verifies the core game mechanics
func TestMakeMove(t *testing.T) {
    game := NewGame("Player1", "Player2")

    // Force a specific die value for predictable testing
    game.CurrentDie = 3

    // Make a move and verify basic state changes
    err := game.MakeMove(0)
    if err != nil {
        t.Errorf("Valid move returned error: %v", err)
    }

    // Verify die was placed in bottom row
    if game.Players[0].Grid[0][0] != 3 {
        t.Errorf("Expected die value 3 in position [0,0], got %d",
            game.Players[0].Grid[0][0])
    }

    // Verify player turn changed
    if game.CurrentPlayer != 1 {
        t.Error("Player turn did not change after move")
    }

    // Verify new die was rolled
    if game.CurrentDie < 1 || game.CurrentDie > 6 {
        t.Error("New die not rolled after move")
    }
}

// TestMatchingDiceRemoval verifies that matching dice are removed from opponent's column
func TestMatchingDiceRemoval(t *testing.T) {
    game := NewGame("Player1", "Player2")

    // Set up scenario: opponent has a 3 in column 0
    game.Players[1].Grid[0][0] = 3

    // Player 0 places a 3 in column 0
    game.CurrentDie = 3
    err := game.MakeMove(0)
    if err != nil {
        t.Fatalf("Valid move returned error: %v", err)
    }

    // Verify opponent's matching die was removed
    if game.Players[1].Grid[0][0] != 0 {
        t.Error("Matching die in opponent's column was not removed")
    }
}

// TestScoring verifies the scoring mechanics
func TestScoring(t *testing.T) {
    game := NewGame("Player1", "Player2")

    // Test single die scoring
    game.Players[0].Grid[0][0] = 3
    game.updateScores()
    if game.Players[0].Score != 3 {
        t.Errorf("Expected score of 3 for single die, got %d", game.Players[0].Score)
    }

    // Test double dice scoring (3 × 2²)
    game.Players[0].Grid[1][0] = 3
    game.updateScores()
    if game.Players[0].Score != 12 {
        t.Errorf("Expected score of 12 for double dice, got %d", game.Players[0].Score)
    }

    // Test triple dice scoring (3 × 3²)
    game.Players[0].Grid[2][0] = 3
    game.updateScores()
    if game.Players[0].Score != 27 {
        t.Errorf("Expected score of 27 for triple dice, got %d", game.Players[0].Score)
    }
}

// TestGameOver verifies game end conditions
func TestGameOver(t *testing.T) {
    game := NewGame("Player1", "Player2")

    // Fill player 1's grid
    for row := 0; row < 3; row++ {
        for col := 0; col < 3; col++ {
            game.Players[0].Grid[row][col] = 1
        }
    }

    // Verify game is over when a grid is full
    if !game.isGameOver() {
        t.Error("Game should be over when a player's grid is full")
    }

    // Test winner determination with different scores
    game.Players[0].Score = 100
    game.Players[1].Score = 90
    game.determineWinner()
    if game.Winner != 0 {
        t.Error("Player 0 should win with higher score")
    }
}

// TestInvalidMoves verifies error handling for invalid moves
func TestInvalidMoves(t *testing.T) {
    game := NewGame("Player1", "Player2")

    // Test move in invalid column
    err := game.MakeMove(-1)
    if err == nil {
        t.Error("Expected error for move in invalid column")
    }

    // Test move in full column
    for row := 0; row < 3; row++ {
        game.Players[0].Grid[row][0] = 1
    }
    err = game.MakeMove(0)
    if err == nil {
        t.Error("Expected error for move in full column")
    }

    // Test move after game is over
    game.GameOver = true
    err = game.MakeMove(1)
    if err == nil {
        t.Error("Expected error for move after game over")
    }
}


// TestScoringCombinations uses table-driven tests to verify various scoring scenarios
func TestScoringCombinations(t *testing.T) {
    // Define our test cases as a slice of anonymous structs
    // Each case describes a specific scoring scenario
    tests := []struct {
        name           string    // descriptive name of the test case
        diceValues     []int     // dice to place in a column
        expectedScore  int       // expected score for this combination
    }{
        {
            name:          "Single die value 6",
            diceValues:    []int{6},
            expectedScore: 6,    // 6 × 1²
        },
        {
            name:          "Two matching dice value 4",
            diceValues:    []int{4, 4},
            expectedScore: 16,   // 4 × 2²
        },
        {
            name:          "Three matching dice value 2",
            diceValues:    []int{2, 2, 2},
            expectedScore: 18,   // 2 × 3²
        },
        {
            name:          "Mixed values in column",
            diceValues:    []int{1, 2, 3},
            expectedScore: 6,    // (1 × 1²) + (2 × 1²) + (3 × 1²)
        },
        {
            name:          "Two pairs in different columns",
            diceValues:    []int{3, 3, 5, 5},
            expectedScore: 32,   // (3 × 2²) + (5 × 2²)
        },
    }

    // Iterate through each test case
    for _, tt := range tests {
        // Use t.Run to create a sub-test for each case
        t.Run(tt.name, func(t *testing.T) {
            game := NewGame("Player1", "Player2")

            // Place dice according to test case
            col := 0
            for i, value := range tt.diceValues {
                game.Players[0].Grid[i][col] = value
            }

            // Calculate score
            game.updateScores()

            // Verify the score matches our expectation
            if game.Players[0].Score != tt.expectedScore {
                t.Errorf("Expected score %d, got %d for case: %s",
                    tt.expectedScore, game.Players[0].Score, tt.name)
            }
        })
    }
}

// TestInvalidMoveScenarios demonstrates table-driven tests for error cases
func TestInvalidMoveScenarios(t *testing.T) {
    tests := []struct {
        name          string
        setupGame     func(*GameState)  // function to set up the test scenario
        moveColumn    int
        expectError   bool
        errorMessage  string
    }{
        {
            name: "Move in negative column",
            setupGame: func(g *GameState) {},  // no setup needed
            moveColumn: -1,
            expectError: true,
            errorMessage: "invalid move",
        },
        {
            name: "Move in column beyond board",
            setupGame: func(g *GameState) {},
            moveColumn: 3,
            expectError: true,
            errorMessage: "invalid move",
        },
        {
            name: "Move in full column",
            setupGame: func(g *GameState) {
                // Fill column 0
                for row := 0; row < 3; row++ {
                    g.Players[0].Grid[row][0] = 1
                }
            },
            moveColumn: 0,
            expectError: true,
            errorMessage: "invalid move",
        },
        {
            name: "Move after game is over",
            setupGame: func(g *GameState) {
                g.GameOver = true
            },
            moveColumn: 0,
            expectError: true,
            errorMessage: "game is already over",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create fresh game state for each test
            game := NewGame("Player1", "Player2")

            // Apply test-specific setup
            tt.setupGame(game)

            // Attempt the move
            err := game.MakeMove(tt.moveColumn)

            // Verify error behavior
            if tt.expectError {
                if err == nil {
                    t.Errorf("Expected error but got none for case: %s", tt.name)
                } else if err.Error() != tt.errorMessage {
                    t.Errorf("Expected error message '%s', got '%s' for case: %s",
                        tt.errorMessage, err.Error(), tt.name)
                }
            } else if err != nil {
                t.Errorf("Unexpected error for case %s: %v", tt.name, err)
            }
        })
    }
}
