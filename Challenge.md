# Go-Knucklebones Autoscaling Challenge

## Background

You've inherited a Go-based implementation of the game Knucklebones, which provides a REST API for creating and playing games. Currently, each server instance can only handle **one game with two players** due to resource constraints. Your challenge is to implement an autoscaling solution that can handle a varying number of players who want to play simultaneously.

## The Problem

The current architecture has the following limitations:

1. Each game server container can only host one game (two players)
2. As player demand increases, we need to automatically spin up new game servers
3. When games are completed, we should shut down the containers to save resources

## Your Task

Implement an autoscaling system using Docker that:

1. Starts with a single game server instance
2. Monitors incoming game creation requests
3. Automatically creates new game server containers when all existing servers are at capacity
4. Intelligently routes players to available servers
5. Gracefully shuts down containers when their games are complete

## Technical Requirements

### Components to Build:

1. **Load Balancer/Orchestrator**: Create a service that:
   - Receives all incoming requests
   - Tracks which game servers are running and their capacity
   - Routes requests to appropriate game servers
   - Spins up new game server containers when needed
   - Shuts down containers when games are complete

2. **Container Management**:
   - Use Docker to manage game server containers
   - Each container should run exactly one instance of the game server
   - Implement proper communication between the orchestrator and containers

3. **Service Discovery**:
   - Implement a way for the load balancer to discover available game servers
   - Track which games are running on which servers

## Testing Your Solution

We've provided a K6 load testing script (`load-test.js`) that simulates an increasing number of players joining and playing games. The script will:

1. Start with a few players
2. Gradually increase the number of concurrent games
3. Complete some games while starting new ones
4. Verify that the system can scale up and down appropriately

Run the test with:

```bash
k6 run load-test.js
```

## Evaluation Criteria

Your solution will be evaluated based on:

1. **Functionality**: Does it correctly handle scaling up and down?
2. **Efficiency**: How effectively does it use resources?
3. **Reliability**: Does it maintain game integrity during scaling operations?
4. **Code Quality**: Is your solution well-organized and maintainable?
5. **Documentation**: Have you clearly explained your approach?

## Getting Started

1. Study the existing codebase to understand how games are created and managed
2. Run the load test to see the current system's limitations
3. Design your autoscaling architecture
4. Implement your solution
5. Test and refine using the provided load testing script

## Submission

Provide the following:

1. Your autoscaling implementation code
2. A `docker-compose.yml` file that configures your solution
3. A README explaining your approach, challenges faced, and how to run your solution
4. Any additional scripts or configurations needed to deploy your solution

Good luck!
