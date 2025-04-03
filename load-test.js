import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter, Rate, Trend } from 'k6/metrics';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

// Custom metrics
const gameCreationErrors = new Counter('game_creation_errors');
const moveErrors = new Counter('move_errors');
const completedGames = new Counter('completed_games');
const createdGames = new Counter('created_games');
const serverResponseTime = new Trend('server_response_time');
const activeGames = new Counter('active_games');

// Base URL
const BASE_URL = 'http://localhost:8080';

// Test stages for ramping up/down
export const options = {
  stages: [
    { duration: '30s', target: 2 }, // Start with 2 games (require 1 container)
    { duration: '30s', target: 5 }, // Ramp up to 5 games (should trigger scaling)
    { duration: '60s', target: 10 }, // Continue to 10 games
    { duration: '30s', target: 20 }, // Stress test with 20 games
    { duration: '60s', target: 10 }, // Scale down to 10
    { duration: '30s', target: 5 }, // Scale down further
    { duration: '30s', target: 0 }, // Complete all games
  ],
  thresholds: {
    'server_response_time': ['p(95)<500'], // 95% of responses should be under 500ms
    'game_creation_errors': ['count<10'], // Fewer than 10 game creation errors
  },
};

// Game state store (simulating state that would be kept in a real client)
const gameStore = {
  games: new Map(),
};

// Create a new game
function createGame() {
  const player1 = `Player_${randomString(5)}`;
  const player2 = `Player_${randomString(5)}`;
  
  const payload = JSON.stringify({
    player1_name: player1,
    player2_name: player2
  });
  
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };
  
  const startTime = new Date();
  const response = http.post(`${BASE_URL}/game`, payload, params);
  const duration = new Date() - startTime;
  serverResponseTime.add(duration);
  
  if (response.status === 200) {
    const gameData = JSON.parse(response.body);
    gameStore.games.set(gameData.id, {
      id: gameData.id,
      state: gameData.state,
      moves: 0,
      movesTarget: Math.floor(Math.random() * 6) + 5, // Random number of moves (5-10)
    });
    createdGames.add(1);
    activeGames.add(1);
    console.log(`Created game ${gameData.id}`);
    return gameData.id;
  } else {
    gameCreationErrors.add(1);
    console.error(`Failed to create game: ${response.status} ${response.body}`);
    return null;
  }
}

// Make a move in an existing game
function makeMove(gameId) {
  if (!gameStore.games.has(gameId)) {
    return false;
  }
  
  const game = gameStore.games.get(gameId);
  const column = Math.floor(Math.random() * 3); // Random column (0-2)
  
  const payload = JSON.stringify({
    column: column
  });
  
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };
  
  const startTime = new Date();
  const response = http.post(`${BASE_URL}/games/${gameId}/move`, payload, params);
  const duration = new Date() - startTime;
  serverResponseTime.add(duration);
  
  if (response.status === 200) {
    const gameData = JSON.parse(response.body);
    game.state = gameData.state;
    game.moves++;
    
    // Check if the game is over or if we've reached our target moves
    if (gameData.state.GameOver || game.moves >= game.movesTarget) {
      completeGame(gameId);
      return true;
    }
    return true;
  } else {
    moveErrors.add(1);
    console.error(`Failed to make move in game ${gameId}: ${response.status} ${response.body}`);
    
    // If we get a 404, the game may have been removed, so we should clean up
    if (response.status === 404) {
      gameStore.games.delete(gameId);
    }
    return false;
  }
}

// Complete a game and notify server
function completeGame(gameId) {
  if (!gameStore.games.has(gameId)) {
    return;
  }
  
  // Notify the server that the game is complete (for cleanup)
  http.post(`${BASE_URL}/games/${gameId}/complete`, null, {
    headers: {
      'Content-Type': 'application/json',
    },
  });
  
  console.log(`Completed game ${gameId} after ${gameStore.games.get(gameId).moves} moves`);
  gameStore.games.delete(gameId);
  completedGames.add(1);
  activeGames.add(-1);
}

// Main function that runs on each VU
export default function() {
  // Get list of active games this VU is playing
  const myActiveGames = Array.from(gameStore.games.keys());
  
  // If we have active games, make moves in them
  if (myActiveGames.length > 0) {
    const gameId = myActiveGames[Math.floor(Math.random() * myActiveGames.length)];
    makeMove(gameId);
  } 
  // Otherwise, try to create a new game
  else {
    createGame();
  }
  
  // Brief pause between actions
  sleep(1);
}
