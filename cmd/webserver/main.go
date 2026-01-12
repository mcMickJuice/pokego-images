package main

import (
	"log"
	"mcmickjuice/pokego/internal/webserver"
)

func main() {
	addr := ":8000"
	webserver := webserver.NewPokemonWebServer(addr)
	if err := webserver.Start(); err != nil {
		log.Fatalf("error starting webserver: %v", err)
	}
}
