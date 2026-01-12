package main

import (
	"fmt"
	"log"
	"mcmickjuice/pokego/internal/webserver"
)

func main() {
	fmt.Println("starting webserver")
	if err := webserver.CreateWebserver(); err != nil {
		log.Fatalf("error starting webserver: %v", err)
	}
}
