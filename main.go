package main

import (
	"knucklebones/server"
	"log"
)

func main() {
	gameServer := server.NewGameServer()
	restAdapter := server.NewRESTAdapter(gameServer, 8080)

	log.Fatal(restAdapter.StartServer())
}
