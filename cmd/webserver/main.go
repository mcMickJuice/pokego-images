package main

import (
	"log"
	"mcmickjuice/pokego/internal/webserver"
)

func main() {
	addr := ":8000"
	server := webserver.NewPokemonWebServer(addr)
	if err := server.Start(); err != nil {
		log.Fatalf("error starting webserver: %v", err)
	}
}
